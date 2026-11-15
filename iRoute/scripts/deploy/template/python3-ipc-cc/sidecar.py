import time
import json
from queue import Queue
import logging
import threading

from config import Config
from message_types import *
from dtc import DTC, SSocket, Poller
from function.handler import workflow_name, writer_up, reader_up, get_lat_dist
from zookeeper import ZookeeperClient

logging.basicConfig(level=logging.INFO)

class Sidecar:
    def __init__(self):
        # record number of func instances
        self.func_info = {}
        self.stage = {}
        self.index = {}
        self.queue = Queue()
        self.socket_queue = Queue()
        self.poller = Poller()
        self.num_dtc = 0
        
        self.get_config()
        self.get_socket()
        self.get_zk()
         
        # func_id: key of self.poller.map
        self.map = {}

        # id for view update, +1 when handle_view_handle is invoked
        self.update_id = "0"
        # update_id: count of instances need to notify when view is updated 
        self.panic_count = {}

        # update_id: func_id, record func_id of instances that need to be deleted after panic
        self.panic_del_ids = {}

        # update_id: func_ids of instances
        self.panic_ids = {}

        # store res of workflow reqs for latency calculation
        self.lats = []

    def get_config(self):
        self.config = Config()
        funcs = self.config.get_funcs()
        for func_name in funcs:
            self.func_info[func_name] = 0
            self.get_index(func_name)

    def get_index(self, func_name):
        for stage_index, stage_fns in enumerate(self.config.config["dag"]):
            try:
                my_func_index = stage_fns.index(func_name)
                my_stage_index = stage_index
                break
            except ValueError as err:
                pass
        self.stage[func_name] = my_stage_index
        self.index[func_name] = my_func_index

    # socket for handshake
    def get_socket(self):
        self.socket = SSocket("0.0.0.0:6000")

    # initialize zookeeper client
    # create nodes for each function instance
    # with write and read ACLs
    def get_zk(self):
        zk_hosts = "zookeeper-headless.default.svc.cluster.local:2181"
        self.zk = ZookeeperClient(
            hosts=zk_hosts,
            username=writer_up[0],
            password=writer_up[1]
        )
        write_acl = self.zk.make_write_acl(writer_up[0], writer_up[1])
        read_acl = self.zk.make_read_acl(reader_up[0], reader_up[1])
        funcs = self.config.get_funcs()
        for func_name in funcs:
            self.zk.create_node(f"/{workflow_name}/{func_name}", data='[]', acl=[write_acl, read_acl])

    def get_fan_in(self, func_name):
        if self.stage[func_name] < len(self.config.config["dag"]) - 1 and len(self.config.config["dag"][self.stage[func_name]]) > 1:
            return [0]
        else:
            return []

    def handle_scale_up(self, message):
        func_name = message.req
        ins_id = self.func_info[func_name] - 1
        if ins_id == 0:
            return

        func_id = f"{func_name}-{ins_id}"
        self.panic_count[self.update_id] = 0

        notify_body = {
            "workflow_name": workflow_name,
            "func_id": func_id,
            "grt_address": "zookeeper-headless.default.svc.cluster.local:2181",
            "acl": reader_up
        }

        req = Message(type_id=MESSAGE_CLIENT_UP, req_id=self.update_id, 
                req=notify_body, index=self.index[func_name])
        if self.stage[func_name] > 0:
            for stage_fn in self.config.config["dag"][self.stage[func_name] - 1]:
                for i in range(self.func_info[stage_fn]):
                    if f"{stage_fn}-{i}" in self.map:
                        self.poller.map[self.map[f"{stage_fn}-{i}"]].send(req)
                        self.panic_ids.append(f"{stage_fn}-{i}")
                        self.panic_count[self.update_id] += 1

        req = Message(type_id=MESSAGE_SERVER_UP, req_id=self.update_id, 
                req=notify_body, index=self.index[func_name])
        if self.stage[func_name] < len(self.config.config["dag"]) - 1:
            for stage_fn in self.config.config["dag"][self.stage[func_name] + 1]:
                for i in range(self.func_info[stage_fn]):
                    if f"{stage_fn}-{i}" in self.map:
                        self.poller.map[self.map[f"{stage_fn}-{i}"]].send(req)
                        self.panic_ids.append(f"{stage_fn}-{i}")
                        self.panic_count[self.update_id] += 1

        self.update_id = str(int(self.update_id) + 1)

    def handle_panic_over(self, message):
        if message.req_id in self.panic_count:
            self.panic_count[message.req_id] -= 1
            # logging.info(f"{self.panic_count[message.req_id]}")

            if self.panic_count[message.req_id] == 0:
                logging.info(f"{message.req_id} panic over")
                self.panic_count.pop(message.req_id, None)

                for func_id in self.panic_ids[message.req_id]:
                    req = Message(type_id=MESSAGE_PANIC_OVER, req_id=None, req=None, index=None)
                    self.poller.map[self.map[func_id]].send(req)

                self.panic_ids.pop(self.message.req_id, None)

                if message.req_id in self.panic_del_ids:
                    func_id = self.panic_del_ids[message.req_id]
                    self.panic_del_ids.pop(message.req_id, None)
                    self.map.pop(func_id, None)

    # handshake with func
    def handshake(self, message, dtc):
        func_name = message.req
        ins_id = self.func_info[func_name]
        func_id = f"{func_name}-{ins_id}"

        func_instances = json.loads(self.zk.get_node(f"/{workflow_name}/{func_name}"))
        func_instances.append(func_id)
        self.zk.set_node(f"/{workflow_name}/{func_name}", json.dumps(func_instances))

        req = {
            "workflow_name": workflow_name,
            "func_id": func_id,
            "cc": (f"cc_{func_id}", f"{func_id}_cc"),
            "grt_address": "zookeeper-headless.default.svc.cluster.local:2181",
            "acl": reader_up
        }

        dtc = DTC(req["cc"][1], req["cc"][0])
        self.poller.add_peer(dtc)
        self.map[func_id] = dtc.rf
        self.num_dtc += 1

        self.func_info[func_name] = self.func_info[func_name] + 1

        response = Message(type_id=MESSAGE_HANDSHAKE, req_id=None, req=req, index=self.index[func_name])

        logging.info(f"Handshake with {func_name}-{ins_id}")
        return response

    def handle_entry(self, func_name):
        req = []
        for i in range(self.func_info[func_name]):
            func_id  = f"{func_name}-{i}"
            if func_id in self.map:
                req.append(func_id)
        response = Message(type_id=MESSAGE_ENTRY, req_id=None, req=req, index=self.index[func_name])
        return response
    
    def handle_scale_down(self, message, dtc):
        func_id = message.req
        func_name = "-".join(func_id.split("-")[:-1])

        func_ids = json.loads(self.zk.get_node(f"/{workflow_name}/{func_name}"))
        try:
            func_ids.remove(func_id)
        except Exception as e:
            pass
        self.zk.set_node(f"/{workflow_name}/{func_name}", json.dumps(func_ids))

        self.panic_count[self.update_id] = 0
        self.panic_ids[self.update_id] = []

        func_stage = self.stage[func_name]

        server_funcs = self.config.config["dag"][func_stage - 1] if func_stage > 0 else []
        req = Message(type_id=MESSAGE_CLIENT_DOWN, req_id=self.update_id, req=func_id, index=self.index[func])
        for func in server_funcs:
            for i in range(self.func_info[func]):
                if f"{func}-{i}" in self.map:
                    self.poller.map[self.map[f"{func}-{i}"]].send(req)
                    self.panic_ids.append(f"{func}-{i}")
                    self.panic_count[self.update_id] += 1

        client_funcs = self.config.config["dag"][func_stage + 1] if func_stage < len(self.config.config["dag"]) - 1 else []
        req = Message(type_id=MESSAGE_SERVER_DOWN, req_id=self.update_id, req=func_id, index=self.index[func])
        for func in client_funcs:
            for i in range(self.func_info[func]):
                if f"{func}-{i}" in self.map:
                    self.poller.map[self.map[f"{func}-{i}"]].send(req)
                    self.panic_ids.append(f"{func}-{i}")
                    self.panic_count[self.update_id] += 1

        self.panic_del_ids[self.update_id] = func_id

        self.update_id = str(int(self.update_id) + 1)

    def handle_closed_dtc(self, dtc):
        dtc_rf = dtc.rf

        func_id = [k for k in self.map if self.map[k] == dtc_rf]
        if func_id:
            func_id = func_id[0]
        else:
            return
        
        req_ids = [k for k, v in self.panic_del_ids.items() if func_id == v]
        if req_ids:
            for req_id in req_ids:
                panic_message = Message(type_id=MESSAGE_PANIC_OVER, req_id=req_id, req=None, index=None)
                self.handle_panic_over(panic_message)
        else:
            func_name = "-".join(func_id.split("-")[:-1])
            message = Message(type_id=MESSAGE_CLIENT_DOWN, req_id=self.update_id, req=func_id, index=self.index[func_name])
            self.handle_scale_down(message, dtc)           

    def socket_recv(self):
        while True:
            recv_message = self.socket.recv()
            self.queue.put((self.socket, recv_message))
            if recv_message.type_id == MESSAGE_EXIT:
                response = Message(type_id=MESSAGE_EXIT, req_id=None, req=None, index=None)
                self.socket.send(response)
                break
            response = self.socket_queue.get()
            self.socket.send(response)

    def pipe_recv(self):
        while True:
            poll_dtcs = self.poller.poll()
            for poll_dtc in poll_dtcs:
                recv_data_list = poll_dtc.recv()
                for recv_message in recv_data_list:
                    if recv_message.type_id == MESSAGE_EXIT:
                        continue
                    self.queue.put((poll_dtc, recv_message))
                if not recv_data_list:
                    self.queue.put((poll_dtc, None))
            if self.num_dtc < 0:
                break

    def run(self):
        srt = threading.Thread(target=self.socket_recv)
        prt = threading.Thread(target=self.pipe_recv)

        srt.start()
        prt.start()

        logging.info(f"Start Centralized Coordinator for {workflow_name}")

        while True:
            dtc, message = self.queue.get()

            if not message:
                self.handle_closed_dtc(dtc)
                continue

            if message.type_id == MESSAGE_OVER:
                self.lats.append((message, time.time()))
            elif message.type_id == MESSAGE_HANDSHAKE:
                response = self.handshake(message, dtc)
                # logging.info(f"{message.req}: {response.req}")
                self.socket_queue.put(response)
                self.handle_scale_up(message)
            elif message.type_id == MESSAGE_PANIC_OVER:
                self.handle_panic_over(message)
            elif message.type_id == MESSAGE_ENTRY:
                response = self.handle_entry(message.req)
                self.socket_queue.put(response)
            elif message.type_id == MESSAGE_LAT_DIS:
                result = "None"
                if self.lats:
                    result = get_lat_dist(self.lats, self.config.config["dag"])
                response = Message(type_id=MESSAGE_LAT_DIS, req_id=None, req=result, index=None)
                self.socket_queue.put(response)
            elif message.type_id == MESSAGE_CLEAR:
                self.lats = []
                response = Message(type_id=MESSAGE_CLEAR, req_id=None, req=None, index=None)
                self.socket_queue.put(response)
            elif message.type_id == MESSAGE_CLIENT_DOWN:
                self.handle_scale_down(message, dtc)
            elif message.type_id == MESSAGE_EXIT:
                break

        # quit
        funcs = self.config.get_funcs()
        for func_name in funcs:
            self.zk.delete_node(f"/{workflow_name}/{func_name}")
        self.zk.stop()
        logging.info("Zookeeper Client Stop")

        self.num_dtc = -1
        dtcs = self.poller.map.values()
        exit_message = Message(type_id=MESSAGE_EXIT, req_id=None, req=None, index=None)
        for dtc in dtcs:
            try:
                dtc.send(exit_message)
            except Exception as e:
                pass
        
        srt.join()
        prt.join()
        
        self.socket.close()
        self.poller.close()

        logging.info("Centralized Coordinator Close")