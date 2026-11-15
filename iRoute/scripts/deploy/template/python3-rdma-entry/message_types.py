import json

MESSAGE_HANDSHAKE = 0
MESSAGE_DATA = 1
MESSAGE_READY = 2
MESSAGE_ROUTE_ERR = 3
MESSAGE_DEL = 4
MESSAGE_SERVER_UP = 5
MESSAGE_SERVER_DOWN = 6
MESSAGE_CLIENT_UP = 7
MESSAGE_CLIENT_DOWN = 8
MESSAGE_PANIC_OVER = 9
MESSAGE_ENTRY = 10
MESSAGE_OVER = 11
MESSAGE_LAT_DIS = 12
MESSAGE_CLEAR = 13
MESSAGE_EXIT = 14
MESSAGE_HANDSHAKE_RDMA = 15

class Message:
    def __init__(self, type_id=None, req_id=None, req=None, index=0):
        """
        type_id: int
        - MESSAGE_HANDSHAKE: transfer dtc info between func and cc
          - req_id: None
          - req: str, func sends name, then cc returns zookeeper info
          - index: int, index of parallel funcs
        - MESSAGE_DATA: result data from upstream
          - req_id: str
          - req: str
          - index: int, index of parallel funcs
        - MESSAGE_READY: all depends ready from downstream
          - req_id: str
          - req: None
          - index: int, index of func that sent this message
        - MESSAGE_ROUTE_ERR: route error from downstream
          - req_id: str
          - req: None
          - index: int, index of func that has route err
        - MESSAGE_DEL: upstream ask to delete data due to route err
          - req_id: str
          - req: None
          - index: None
        - MESSAGE_SERVER_UP or MESSAGE_CLIENT_UP: scale up from controller, enter panic mode,
          - req_id: str, update_id
          - req: None
          - index: int, index of func needs scale up
        - MESSAGE_SERVER_DOWN or MESSAGE_CLIENT_DOWN: scale down from controller, enter panic mode,
          - req_id, str, update_id
          - req: str, func_id
          - index: int, index of func needs scale down
        - MESSAGE_PANIC_OVER: all funcs have updated the routing table, exit panic mode
          - req_id: str, update id
          - req: None
          - index: None
        - MESSAGE_ENTRY: get the server ips of entry funcs
          - req_id: None
          - req: str, func_name
          - index: None
        - MESSAGE_OVER: funcs of final stage send to cc for calculating latency
          - req_id: str
          - req: str
          - index: int
        - MESSAGE_LAT_DIS: latency distribution from cc
          - req_id: None
          - req: None
          - index: None
        - MESSAGE_CLEAR: clear the stored latencies in cc
          - req_id: None
          - req: None
          - index: None
        - MESSAGE_EXIT: exit the sidecar
          - req_id: None
          - req: None
          - index: None
        - MESSAGE_HANDSHAKE_RDMA: handshake RDMA information
          - req_id: None
          - req: str, '{gid},{qp_num},{recv_mr.buf},{recv_mr.rkey}'
          - index: int (nqp)
        """
        self.type_id = type_id
        self.req_id = req_id
        self.req = req
        self.index = index

    def dumps(self):
        message = {"type": self.type_id}
        if not self.req_id is None:
            message["req_id"] = self.req_id
        if not self.req is None:
            message["req"] = self.req
        if not self.index is None:
            message["index"] = self.index

        return json.dumps(message)

    def loads(self, message_str):
        message = json.loads(message_str)

        self.type_id = message["type"]
        if "req_id" in message:
            self.req_id = message["req_id"]
        if "req" in message:
            self.req = message["req"]
        if "index" in message:
            self.index = message["index"]