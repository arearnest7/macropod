import time
import json
from queue import Queue
import logging
import socket

from config import Config
from message_types import *
from dtc import Server
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
        
        self.get_config()
        self.get_socket()
        self.get_zk()
         
        # func_id: ip, record server_ip of func instance
        self.ips = {}
        # func_id -> socket fd
        self.client_fds = {}

        # id for view update, +1 when handle_view_handle is invoked
        self.update_id = "0"
        # update_id: count of instances need to notify when view is updated 
        self.panic_count = {}
        self.panic_fds = {}

        # update_id: func_id, record func_id of instances that need to be deleted after panic
        self.panic_func_ids = {}

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
        self.socket = Server("0.0.0.0", 6000)

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
            self.zk.create_node(f"/{workflow_name}/{func_name}", data='{}', acl=[write_acl, read_acl])

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
        self.panic_fds[self.update_id] = []

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
                    if f"{stage_fn}-{i}" in self.client_fds:
                        self.socket.send(req, self.socket.map[self.client_fds[f"{stage_fn}-{i}"]])
                        self.panic_fds[self.update_id].append(self.client_fds[f"{stage_fn}-{i}"])
                        self.panic_count[self.update_id] += 1

        self.update_id = str(int(self.update_id) + 1)

    def handle_panic_over(self, message):
        if message.req_id in self.panic_count:
            self.panic_count[message.req_id] -= 1
            # logging.info(f"{self.panic_count[message.req_id]}")

            if self.panic_count[message.req_id] == 0:
                logging.info(f"{message.req_id} panic over")
                self.panic_count.pop(message.req_id, None)

                for fd in self.panic_fds[message.req_id]:
                    req = Message(type_id=MESSAGE_PANIC_OVER, req_id=None, req=None, index=None)
                    self.socket.send(req, self.socket.map[fd])

                self.panic_fds.pop(message.req_id, None)

                if message.req_id in self.panic_func_ids:
                    func_id = self.panic_func_ids[message.req_id]
                    self.ips.pop(func_id, None)
                    self.client_fds.pop(func_id, None)
                    self.panic_func_ids.pop(message.req_id, None)

    # handshake with func
    def handshake(self, message, dtc):
        func_name = message.req
        ins_id = self.func_info[func_name]
        func_id = f"{func_name}-{ins_id}"

        ins_address = self.socket.address[dtc.fileno()]
        ins_host = ins_address.split(":")[0]

        func_instances = json.loads(self.zk.get_node(f"/{workflow_name}/{func_name}"))
        func_instances[func_id] = ins_host
        self.zk.set_node(f"/{workflow_name}/{func_name}", json.dumps(func_instances))

        req = {
            "workflow_name": workflow_name,
            "func_id": func_id,
            "grt_address": "zookeeper-headless.default.svc.cluster.local:2181",
            "acl": reader_up
        }

        self.ips[func_id] = ins_host
        self.client_fds[func_id] = dtc.fileno()

        self.func_info[func_name] = self.func_info[func_name] + 1

        response = Message(type_id=MESSAGE_HANDSHAKE, req_id=None, req=req, index=self.index[func_name])

        logging.info(f"Handshake with {func_name}-{ins_id} at {ins_address}")
        return response

    def handle_entry(self, func_name):
        req = []
        for i in range(self.func_info[func_name]):
            req.append((i, self.ips[f"{func_name}-{i}"]))
        response = Message(type_id=MESSAGE_ENTRY, req_id=None, req=req, index=self.index[func_name])
        return response
    
    def handle_scale_down(self, message, dtc):
        func_id = message.req
        func_name = "-".join(func_id.split("-")[:-1])

        func_instances = json.loads(self.zk.get_node(f"/{workflow_name}/{func_name}"))
        func_instances.pop(func_id, None)
        self.zk.set_node(f"/{workflow_name}/{func_name}", json.dumps(func_instances))

        self.panic_count[self.update_id] = 0
        self.panic_fds[self.update_id] = []

        # just send MESSAGE_CLIENT_DOWN to the upstream funcs of the removed func instance
        func_stage = self.stage[func_name]
        target_funcs = self.config.config["dag"][func_stage - 1] if func_stage > 0 else []
        req = Message(type_id=MESSAGE_CLIENT_DOWN, req_id=self.update_id, req=func_id, index=self.index[func])
        for func in target_funcs:
            for i in range(self.func_info[func]):
                if f"{func}-{i}" in self.client_fds:
                    self.socket.send(req, self.socket.map[self.client_fds[f"{func}-{i}"]])
                    self.panic_fds[self.update_id].append(self.client_fds[f"{func}-{i}"])
                    self.panic_count[self.update_id] += 1

        self.panic_func_ids[self.update_id] = func_id

        self.update_id = str(int(self.update_id) + 1)

    def handle_closed_socket(self, dtc):
        dtc_fd = dtc.fileno()

        func_id = [k for k in self.client_fds if self.client_fds[k] == dtc_fd]
        if func_id:
            func_id = func_id[0]
        else:
            return
        
        req_ids = [k for k, v in self.panic_func_ids.items() if func_id == v]
        if not req_ids:
            return
        
        for req_id in req_ids:
            self.panic_fds[req_id].remove(dtc_fd)
            panic_message = Message(type_id=MESSAGE_PANIC_OVER, req_id=req_id, req=None, index=None)
            self.handle_panic_over(panic_message)
        
        func_name = "-".join(func_id.split("-")[:-1])
        message = Message(type_id=MESSAGE_CLIENT_DOWN, req_id=self.update_id, req=func_id, index=self.index[func_name])
        self.handle_scale_down(message, dtc)

    def run(self):
        logging.info(f"Start Centralized Coordinator for {workflow_name}")
        flag = True
        while flag:
            poll_results = self.socket.poll()
            for dtc, message in poll_results:
                if message is None:
                    self.handle_closed_socket(dtc)
                    continue
                elif isinstance(message, str):
                    req = "Please send message instead of str"
                    response = Message(type_id=MESSAGE_DATA, req_id=None, req=req, index=None)
                    self.socket.send(response, dtc)
                if message.type_id == MESSAGE_OVER:
                    self.lats.append((message, time.time()))
                elif message.type_id == MESSAGE_HANDSHAKE:
                    response = self.handshake(message, dtc)
                    # logging.info(f"{message.req}: {response.req}")
                    self.socket.send(response, dtc)
                    self.handle_scale_up(message)
                elif message.type_id == MESSAGE_PANIC_OVER:
                    self.handle_panic_over(message)
                elif message.type_id == MESSAGE_ENTRY:
                    response = self.handle_entry(message.req)
                    self.socket.send(response, dtc)
                elif message.type_id == MESSAGE_LAT_DIS:
                    result = "None"
                    if self.lats:
                        result = get_lat_dist(self.lats, self.config.config["dag"])
                    response = Message(type_id=MESSAGE_LAT_DIS, req_id=None, req=result, index=None)
                    self.socket.send(response, dtc)
                elif message.type_id == MESSAGE_CLEAR:
                    self.lats = []
                    response = Message(type_id=MESSAGE_CLEAR, req_id=None, req=None, index=None)
                    self.socket.send(response, dtc)
                elif message.type_id == MESSAGE_CLIENT_DOWN:
                    self.handle_scale_down(message, dtc)
                elif message.type_id == MESSAGE_EXIT:
                    response = Message(type_id=MESSAGE_EXIT, req_id=None, req=None, index=None)
                    self.socket.send(response, dtc)
                    flag = False

        # quit
        funcs = self.config.get_funcs()
        for func_name in funcs:
            self.zk.delete_node(f"/{workflow_name}/{func_name}")
        self.zk.stop()
        logging.info("Zookeeper Client Stop")

        dtcs = self.socket.client_socket
        exit_message = Message(type_id=MESSAGE_EXIT, req_id=None, req=None, index=None)
        for dtc in dtcs:
            try:
                self.socket.send(exit_message, dtc)
            except Exception as e:
                pass
        self.socket.clean()
        logging.info("Centralized Coordinator Close")