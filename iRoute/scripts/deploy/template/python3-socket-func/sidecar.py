import json
import logging
import traceback

from dtc import Client, Server
from consistent import HashTable
from message_types import *
from zookeeper import ZookeeperClient

import time

logging.basicConfig(level=logging.INFO)

class Sidecar:
    def __init__(self, handle, config):
        self.func_name = config.func_name
        # logging.info(self.func_name)

        # handshake with cc
        workflow_name, grt_address, read_up = self.handshake(config)

        self.fan_in = config.fan_in
        self.server = Server("0.0.0.0", 6000)
        self.server.add_dtc(self.cc)

        # fd of upstream func -> index of func
        self.servers = {}
        # str(index) -> fds
        self.server_index_list = {}
        for i in range(self.fan_in):
            self.server_index_list[str(i)] = []

        client_funcs = config.get_clients()
        self.fan_out = len(client_funcs)
        
        # get clients from grt
        clients = self.get_clients_from_grt(grt_address, workflow_name, read_up, client_funcs)

        # 2-dim array. 1st are fan-out funcs, 2nd are (func_id, fd)
        # tuple[][]: array of (func_id, fd) for downstream funcs
        self.clients = []
        self.fan_in_clients = config.get_fan_in()
        # [], hash table for fan_in, int for others
        self.hash_tables = []
        # client_id - > ip
        self.client_ips = {}
        for client_index, client_infos in enumerate(clients):
            # client_info: {client_id: client_host}
            self.clients.append([])
            for client_id, client_host in client_infos.items():
                self.client_ips[client_id] = client_host
                client = Client(client_host, 6000)
                self.server.add_dtc(client)
                self.clients[client_index].append((client_id, client.socket.fileno()))
                
                if client_index in self.fan_in_clients:
                    self.hash_tables.append(HashTable())
                    self.hash_tables[-1].add_server(client_id)
                else:
                    self.hash_tables.append(0)

        # worker funcntion
        self.handle = handle
        # index of parallel funcs
        self.index = config.index

        # inputs from upstream funcs
        self.inputs = {}
        # save outputs temporarily
        self.outputs = {}
        # req_id:index that have wrong hash route, used by servers
        self.replays = []
        # req_id:index -> dtc socket, indicate replay client, used by clients
        self.dtc_map = {}

        self.panic = False
    
    # handshake with cc
    def handshake(self, config):
        cc_name = config.cc_name
        cc_ip = f"{cc_name}.openfaas-fn.svc.cluster.local"
        cc_port = 6000
        self.cc = Client(cc_ip, cc_port)

        message = Message(type_id=MESSAGE_HANDSHAKE, req_id=None, req=self.func_name, index=None)
        self.cc.send(message)

        response = self.cc.recv()

        handshake_info = response.req
        self.func_id = handshake_info["func_id"]
        return handshake_info["workflow_name"], handshake_info["grt_address"], handshake_info["acl"]
    
    def get_clients_from_grt(self, grt_address, workflow_name, read_up, func_names):
        zk = ZookeeperClient(grt_address, read_up[0], read_up[1])

        clients = []

        for func_name in func_names:
            while True:
                func_instances = json.loads(zk.get_node(f"/{workflow_name}/{func_name}"))
                if func_instances:
                    clients.append(func_instances)
                    break
                else:
                    time.sleep(1)

        zk.stop()

        return clients

    def choose_client(self, client_index, req_id):
        # only 1 client
        if len(self.clients[client_index]) == 1:
            client_fd = self.clients[client_index][0][1]
        # multiple clients and fan-in, hash
        elif client_index in self.fan_in_clients:
            client_ip = self.client_ips[self.hash_tables[client_index].get_server(req_id)]
            client_fd = self.server.pam[client_ip]
        # multiple clients but non-fan-in, round-robin
        else:
            choose_index = self.hash_tables[client_index]
            self.hash_tables[client_index] = (self.hash_tables[client_index] + 1) % len(self.clients[client_index])
            client_fd = self.clients[client_index][choose_index][1]

        return self.server.map[client_fd]

    # select_dtc is socket object
    def sender(self, select_dtc, message):
        self.server.send(message, select_dtc)

    def worker(self, req_id, req):
        # logging.info("start worker")
        self.n_req_interval += 1
        self.n_req += 1

        exec_start = time.time()

        try:
            res = self.handle(req)
        except Exception as e:
            traceback.print_exc()
            self.n_req -= 1
            return

        # logging.info("finish handle")

        exec_time = time.time() - exec_start

        if self.avg_exec_time:
            avg_exec_time = (self.avg_exec_time * (self.n_req - 1) + exec_time) / self.n_req
            self.avg_exec_time = avg_exec_time
        else:
            self.avg_exec_time = exec_time

        # logging.info(f"{self.func_name}: {res}")

        if self.fan_out > 0:
            # store result until all downstream funcs are ready
            self.outputs[req_id] = {"ready": 0, "req": res}
            message = Message(type_id=MESSAGE_DATA, req_id=req_id, req=res, index=self.index)
            for client_index in range(self.fan_out):
                replay_id = f"{req_id}:{client_index}"
                # sometimes index=0 is not the critical, but still adopt its route
                if replay_id in self.dtc_map:
                    self.sender(self.dtc_map[replay_id], message)
                    self.dtc_map.pop(replay_id, None)
                else:
                    # logging.info(f"start send req to client {client_index}")
                    chosen_client = self.choose_client(client_index, req_id)
                    self.sender(chosen_client, message)
                    self.outputs[req_id][str(client_index)] = chosen_client
        else:
            # pass
            # only for getting latency distribution
            message = Message(type_id=MESSAGE_OVER, req_id=req_id, req=res, index=self.index)
            self.cc.send(message)
            # logging.info(f"{self.func_name}: {res}")

    # massage.type_id == MESSAGE_DATA
    def handle_func_req(self, poll_dtc, message):
        fd = poll_dtc.fileno()
        if fd not in self.servers:
            self.servers[fd] = message.index
            self.server_index_list[str(message.index)].append(fd)

        req_id = message.req_id

        # if self.fan_in == 1:
        #     self.worker(req_id, message.req)
        #     return

        if req_id in self.inputs:
            self.inputs[req_id]["ready"] = self.inputs[req_id]["ready"] + 1
        else:
            self.inputs[req_id] = {"ready": 1, "req": [None] * self.fan_in}

        self.inputs[req_id]["req"][message.index] = message.req

        replay_id = f"{req_id}:{message.index}"
        # complete req
        if self.inputs[req_id]["ready"] == self.fan_in:
            if self.fan_in == 1:
                data = self.inputs[req_id]["req"][0]
            else:
                if isinstance(message.req, str):
                    # merge data for string data
                    data = json.dumps(self.inputs[req_id]["req"])
                else:
                    # merge data for non-string data
                    data = self.inputs[req_id]["req"]

            self.worker(req_id, data)

            # del self.inputs[req_id]
            self.inputs.pop(req_id, None)

            # return

            response = Message(type_id=MESSAGE_READY, req_id=req_id, req=None, index=self.index)

            # if replay, notify all upstream clients
            if replay_id in self.replays:
                for client_fd in self.server_index_list[str(message.index)]:
                    self.sender(self.server.map[client_fd], response)
                self.replays.remove(replay_id)
            # else, only response client
            else:
                self.sender(poll_dtc, response)
        # critical path, ask replay if in panic mode
        elif message.index == 0 and self.panic:
            response = Message(type_id=MESSAGE_ROUTE_ERR, req_id=req_id, req=None, index=self.index)
            # notify all upstream clients
            for index, req in enumerate(self.inputs[req_id]["req"]):
                if req is None:
                    self.replays.append(f"{req_id}:{index}")
                    for client_fd in self.server_index_list[str(message.index)]:
                        self.sender(self.server.map[client_fd], response)
        # reply of replay
        elif replay_id in self.replays:
            response = Message(type_id=MESSAGE_READY, req_id=req_id, req=None, index=self.index)
            for client_fd in self.server_index_list[str(message.index)]:
                self.sender(self.server.map[client_fd], response)
            self.replays.remove(replay_id)
            
    # massage.type_id == MESSAGE_READY
    def handle_depends_ready(self, message):
        req_id = message.req_id
        # check if del outputs
        if req_id in self.outputs:
            self.outputs[req_id]["ready"] = self.outputs[req_id]["ready"] + 1

            if self.outputs[req_id]["ready"] == self.fan_out:
                # del self.outputs[req_id]
                self.outputs.pop(req_id, None)
        else:
            logging.error("---some wrong here?---")

        # check if del replay
        replay_id = f"{req_id}:{message.index}"
        if replay_id in self.replays:
            self.replays.remove(replay_id)
        if replay_id in self.dtc_map:
            # del self.dtc_map[replay_id]
            self.dtc_map.pop(replay_id, None)

    # massage.type_id == MESSAGE_ROUTE_ERR
    def handle_route_error(self, poll_dtc, message):
        req_id = message.req_id
        # route error, replay
        if req_id in self.outputs:
            response = Message(type_id=MESSAGE_DATA, req_id=req_id, req=self.outputs[req_id]["req"], index=self.index)
            self.sender(poll_dtc, response)

            # ask the original client to delete data
            del_message = Message(type_id=MESSAGE_DEL, req_id=req_id, req=None, index=None)
            self.sender(self.outputs[req_id][self.index], del_message)
        # maybe still execute, store dtc for future
        else:
            self.dtc_map[f"{req_id}:{message.index}"] = poll_dtc

    # check if some reqs in stable mode needs replay
    def check_replay(self):
        for req_id in self.inputs:
            response = None
            # notify all upstream clients
            for index, req in enumerate(self.inputs[req_id]["req"]):
                # if critical path has data
                if index == 0 and req is not None:
                    break
                else:
                    response = Message(type_id=MESSAGE_ROUTE_ERR, req_id=req_id, req=None, index=self.index)
                if req is None:
                    self.replays.append(f"{req_id}:{index}")
                    for client_fd in self.server_index_list[str(index)]:
                        self.sender(self.server.map[client_fd], response)

    # message.type_id == MESSAGE_SERVER_UP or MESSAGE_CLIENT_UP
    # only MESSAGE_CLIENT_UP in TCP SOCKET, MESSAGE_SERVER_UP happends in socket connect
    def handle_scale_up(self, poll_dtc, message):
        # add clients
        client_id = (message.req)["func_id"]
        workflow_name = message.req["workflow_name"]
        grt_address = message.req["grt_address"]
        reader_up = message.req["acl"]

        client_name = "-".join(client_id.split("-")[:-1])

        zk = ZookeeperClient(grt_address, reader_up[0], reader_up[1])
        func_instances = json.loads(zk.get_node(f"/{workflow_name}/{client_name}"))
        zk.stop()
        
        client_ip = func_instances[client_id]
        
        self.client_ips[client_id] = client_ip
        
        client = Client(client_ip, 6000)
        self.server.add_dtc(client)
        self.clients[message.index].append((client_id, client.socket.fileno()))

        if message.index in self.fan_in_clients:
            self.hash_tables[message.index].add_server(client_id)

        self.panic = True

        response = Message(type_id=MESSAGE_PANIC_OVER, req_id=message.req_id, req=None, index=self.index)
        self.sender(poll_dtc, response)

        self.check_replay()

    # message.type_id == MESSAGE_SERVER_DOWN or MESSAGE_CLIENT_DOWN
    def handle_scale_down(self, poll_dtc, message):
        # del servers; message from upstream func that will exit
        if message.type_id == MESSAGE_SERVER_DOWN:
            self.servers.pop(poll_dtc.fileno(), None)
            self.server_index_list[str(message.index)].remove(poll_dtc.fileno())
            # self.server.del_dtc(poll_dtc)            
        # del clients
        else:
            client_id = message.req
            client_ip = self.client_ips[client_id]
            client_fd = self.server.pam[client_ip]
            self.clients[message.index].remove((client_id, client_fd))
            # self.server.del_dtc(self.server.map[client_fd])

            if message.index in self.fan_in_clients:
                self.hash_tables[message.index].del_server(client_id)
                
            self.panic = True

            response = Message(type_id=MESSAGE_PANIC_OVER, req_id=message.req_id, req=None, index=self.index)
            self.sender(poll_dtc, response)

            self.check_replay()

    def handle_closed_socket(self, poll_dtc):
        dtc_fd = poll_dtc.fileno()
        if dtc_fd in self.servers:
            self.servers.pop(dtc_fd, None)
            indexs = [k for k, v in self.server_index_list.items() if dtc_fd in v]
            for index in indexs:
                self.server_index_list[index].remove(dtc_fd)

    def update_arrival_rate(self):
        self.arrival_rate = self.n_req_interval * 1.0 / (time.time() - self.start_exec_time - 5)
        
        if self.arrival_rate > 0:
            if self.arrival_rate > self.n_req_interval:
                logging.info(f"{self.func_name}: recv {self.n_req_interval} in {time.time() - self.start_exec_time - 5}s")
            else:
                logging.info(f"{self.func_name}: RPS={self.arrival_rate}/s")

            if self.avg_exec_time:
                logging.info(f"{self.func_name}: QPS={1.0 / self.avg_exec_time}/s")

        self.n_req_interval = 0
        self.start_exec_time = time.time()

    def run(self):
        self.arrival_rate = None
        self.start_exec_time = time.time()
        self.n_req = 0
        self.n_req_interval = 0
        self.avg_exec_time = None

        logging.info(f"{self.func_name} LC started")

        flag = True
        while flag:
            poll_results = self.server.poll()
            if not poll_results:
                # no message, update arrival rate
                self.update_arrival_rate()
                continue
            for poll_dtc, message in poll_results:
                if message is None:
                    self.handle_closed_socket(poll_dtc)
                    continue

                if message.type_id == MESSAGE_DATA:
                    # logging.info("Start handling req")
                    self.handle_func_req(poll_dtc, message)
                elif message.type_id == MESSAGE_READY:
                    self.handle_depends_ready(message)
                elif message.type_id == MESSAGE_ROUTE_ERR:
                    self.handle_route_error(poll_dtc, message)
                elif message.type_id == MESSAGE_DEL:
                    self.inputs.pop(message.req_id, None)
                elif message.type_id == MESSAGE_SERVER_UP or message.type_id == MESSAGE_CLIENT_UP:
                    self.handle_scale_up(poll_dtc, message)
                elif message.type_id == MESSAGE_SERVER_DOWN or message.type_id == MESSAGE_CLIENT_DOWN:
                    self.handle_scale_down(poll_dtc, message)
                elif message.type_id == MESSAGE_PANIC_OVER:
                    self.panic = False
                    logging.info(f"{self.func_name} panic over")
                elif message.type_id == MESSAGE_EXIT:
                    flag = False
                    continue

        self.quit()

        logging.info(f"{self.func_name} LC closed")

    def quit(self):
        req = Message(type_id=MESSAGE_CLIENT_DOWN, req_id=None, req=self.func_id, index=self.index)
        self.cc.send(req)

        req = Message(type_id=MESSAGE_SERVER_DOWN, req_id=None, req=None, index=self.index)
        if self.fan_out:
            for client_info in self.clients:
                for client_id, client_fd in client_info:
                    self.server.send(req, self.server.map[client_fd])

        self.server.clean()