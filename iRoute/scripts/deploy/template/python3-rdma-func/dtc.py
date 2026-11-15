import socket
import select
import time
import pickle
import struct
import math

from pyverbs.addr import AH, AHAttr, GlobalRoute, GID
from pyverbs.cq import CQ, CompChannel
from pyverbs.device import Context
from pyverbs.enums import *
from pyverbs.mr import MR
from pyverbs.pd import PD
from pyverbs.qp import QP, QPCap, QPInitAttr, QPAttr
from pyverbs.wr import SGE, RecvWR, SendWR

class RDMA:
    def __init__(self, func_name, args):
        self.func_name = func_name

        self.args = {
            'ib_dev': 'mlx5_1', 'sg_depth': 1, 'inline_size': 0, 'mtu': 2, 'rx_depth': 2500, 
            'size': 1024, 'sl': 0, 'operation_type': IBV_WR_RDMA_WRITE_WITH_IMM , 'tx_depth': 2500, 
            'qp_type': IBV_QPT_RC, 'gid_index': 3
        }
        self.args.update(args)

        # number of sgl for each mr
        self.N = 50
        if self.args.get("num_sgl", None):
            self.N = self.args["num_sgl"]

        self.empty = " " * self.args["size"]
        # k: IDENTITY, v: qp
        self.qps = {}
        self.nqp = 0
        # if mr for recv - k: IDENTITY + '-1', v: {"mr": mr, "sgls": sgls, "index": index of sgl currently in use}
        # if mr for send - k: IDENTITY + '-0', v: {"mr": mr, "sgl": sgl}
        self.mrs = {}
        # wr_id -> mr_id
        self.map = {}
        # mr_id -> wr_id
        self.pam = {}

        # nqp -> local_info
        self.handshake_info = {}

        # mr_id -> {"chunks": {index: chunk}, "start": int, "num_chunks": int}
        # store chunks of big data
        self.bk = {}

        self.rdma_init()

    def rdma_init(self):
        ctx = Context(name=self.args['ib_dev'])
        self.pd = PD(ctx)
        self.comp_ch = CompChannel(ctx)
        self.cq = CQ(ctx, 100, None, self.comp_ch)
        self.cq.req_notify()

        cap = QPCap(max_send_wr=self.args['tx_depth'], max_recv_wr=self.args['rx_depth'], max_send_sge=self.args['sg_depth'],
            max_recv_sge=self.args['sg_depth'], max_inline_data=self.args['inline_size'])
        self.qp_init_attr = QPInitAttr(qp_type=self.args['qp_type'], scq=self.cq, rcq=self.cq, cap=cap, sq_sig_all=True)

        self.gid = ctx.query_gid(port_num=1, index=self.args['gid_index'])

    # recv message
    def mr_read(self, mr_id, index):
        mr = self.mrs[mr_id]["mr"]
        buf = index * self.args['size']

        recv_res = None
        
        if mr_id in self.bk:
            first_buf = mr.read(4, buf)
            chunk_len = struct.unpack("!I", first_buf[:4])[0]

            chunk_data = mr.read(chunk_len, buf + 4)
            self.bk[mr_id]["chunks"][str(index)] = chunk_data

            if len(self.bk[mr_id]["chunks"]) == self.bk[mr_id]["num_chunks"]:
                indexs = [(self.bk[mr_id]["start"] + i) % self.N for i in range(self.bk[mr_id]["num_chunks"])]
                recv_res = b''.join([self.bk[mr_id]["chunks"][str(i)] for i in indexs])
                self.bk.pop(mr_id, None)
        else:
            first_buf = mr.read(8, buf)
            total_len = struct.unpack("!I", first_buf[:4])[0]
            chunk_len = struct.unpack("!I", first_buf[4:8])[0]

            first_payload_cap = self.args['size'] - 8
            other_payload_cap = self.args['size'] - 4

            remaining = total_len - chunk_len
            num_other = math.ceil(remaining / other_payload_cap) if remaining > 0 else 0
            num_chunks = 1 + num_other

            chunk_data = mr.read(chunk_len, buf + 8)

            if num_chunks > 1:
                self.bk[mr_id] = {"chunks": {str(index): chunk_data}, "start": index, "num_chunks": num_chunks}
            else:
                recv_res = chunk_data
        
        self.add_wr(self.qps[mr_id[:-2]], self.pam[mr_id], self.mrs[mr_id]["sgl"][index])

        if recv_res:
            return pickle.loads(recv_res)
        else:
            return None


    # send message
    def mr_write(self, mr_id, send_obj):
        data_b = pickle.dumps(send_obj)

        data_len = len(data_b)

        first_payload_cap = self.args['size'] - 8
        other_payload_cap = self.args['size'] - 4

        first_taken = min(first_payload_cap, data_len)
        remaining = data_len - first_taken
        num_other = math.ceil(remaining / other_payload_cap) if remaining > 0 else 0
        num_chunks = 1 + num_other

        offset = 0
        for i in range(num_chunks):
            if i == 0:
                cap = first_payload_cap
            else:
                cap = other_payload_cap

            chunk_data = data_b[offset: offset + cap]
            chunk_len = len(chunk_data)

            if i == 0:
                # first chunk
                header = struct.pack("!I", data_len) + struct.pack("!I", chunk_len)
                self.mrs[mr_id]["mr"].write(header, len(header), 0)
                self.mrs[mr_id]["mr"].write(chunk_data, chunk_len, 8)
            else:
                # other chunks
                header = struct.pack("!I", chunk_len)
                self.mrs[mr_id]["mr"].write(header, len(header), 0)
                self.mrs[mr_id]["mr"].write(chunk_data, chunk_len, 4)

            offset += chunk_len

            index = self.mrs[mr_id]["index"]

            wr = SendWR(self.pam[mr_id], opcode=self.args['operation_type'], num_sge=1, 
                    imm_data=index, sg=self.mrs[mr_id]["sgl"])
            self.mrs[mr_id]["index"] += 1
            if self.mrs[mr_id]["index"] == self.N:
                self.mrs[mr_id]["index"] = 0
            wr.set_wr_rdma(self.mrs[mr_id]["rkey"], self.mrs[mr_id]["addr"] + index * self.args['size'])

            self.cq.req_notify()
            self.qps[mr_id[:-2]].post_send(wr)
            self.cq.req_notify()

    # direct send
    def send(self, mr_id, send_obj):
        self.mr_write(mr_id, send_obj)

    def poll(self):
        self.comp_ch.get_cq_event(self.cq)
        self.cq.req_notify()
        self.cq.ack_events(1)

        poll_results = []
        while True:
            wc_num, wc_list = self.cq.poll()

            if wc_num > 0:
                if wc_list[0].opcode == int(IBV_WC_RECV_RDMA_WITH_IMM):
                    mr_id = self.map[str(wc_list[0].wr_id)]
                    index = wc_list[0].imm_data
                    recv_data = self.mr_read(mr_id, index)
                    if recv_data is not None:
                        poll_results.append((mr_id, recv_data))
                continue
            else:
                break

        return poll_results

    def prepare_dtc(self, index):
        qp = QP(self.pd, self.qp_init_attr)
        send_mr = MR(self.pd, self.args['size'], IBV_ACCESS_LOCAL_WRITE | IBV_ACCESS_REMOTE_WRITE | IBV_ACCESS_REMOTE_READ)
        recv_mr = MR(self.pd, self.args['size'] * self.N, IBV_ACCESS_LOCAL_WRITE | IBV_ACCESS_REMOTE_WRITE | IBV_ACCESS_REMOTE_READ)

        self.handshake_info[str(index)] = (qp, send_mr, recv_mr)

        local_info = f"{str(self.gid)},{str(qp.qp_num)},{str(recv_mr.buf)},{str(recv_mr.rkey)}"

        return local_info

    def add_dtc(self, remote_info_str, IDENTITY, index):
        gid_str, qpn_str, remote_addr, remote_key = remote_info_str.split(",")

        remote_info = {"gid": GID(gid_str), "qpn": int(qpn_str), "addr": int(remote_addr), "rkey": int(remote_key)}

        qp, send_mr, recv_mr = self.handshake_info[str(index)]

        gr = GlobalRoute(dgid=remote_info['gid'], sgid_index=self.args['gid_index'])
        ah_attr = AHAttr(gr=gr, is_global=1, port_num=1)

        qa = QPAttr()
        qa.ah_attr = ah_attr
        qa.dest_qp_num = remote_info['qpn']
        qa.path_mtu = self.args['mtu']
        qa.max_rd_atomic = 1
        qa.max_dest_rd_atomic = 1
        qa.qp_access_flags = IBV_ACCESS_REMOTE_WRITE | IBV_ACCESS_REMOTE_READ | IBV_ACCESS_LOCAL_WRITE
        qp.to_rts(qa)

        self.qps[IDENTITY] = qp
        smr_id = f"{IDENTITY}-0"
        rmr_id = f"{IDENTITY}-1"

        self.map[str(index * 2)] = smr_id
        self.map[str(index * 2 + 1)] = rmr_id
        self.pam[smr_id] = index * 2
        self.pam[rmr_id] = index * 2 + 1

        send_sgl = self.add_sgls(qp, send_mr, None)
        recv_sgls = self.add_sgls(qp, recv_mr, index * 2 + 1)

        self.mrs[smr_id] = {"mr": send_mr, "sgl": send_sgl, "index": 0, "rkey": remote_info["rkey"], "addr": remote_info["addr"]}
        self.mrs[rmr_id] = {"mr": recv_mr, "sgl": recv_sgls}

        del self.handshake_info[str(index)]

    def add_sgls(self, qp, mr, wr_id):
        if wr_id:
            sgls = []
            for i in range(self.N):
                sgls.append([SGE(mr.buf + i * self.args['size'], self.args['size'], mr.lkey)])
            self.add_wrs(qp, sgls, wr_id)
            return sgls
        else:
            sgl = [SGE(mr.buf, self.args['size'], mr.lkey)]
            return sgl

    def add_wrs(self, qp, sgls, wr_id):
        for sgl in sgls:
            wr = RecvWR(wr_id, len(sgl), sgl)
            self.cq.req_notify()
            qp.post_recv(wr)

    def add_wr(self, qp, wr_id, sgl):
        wr = RecvWR(wr_id, len(sgl), sgl)
        self.cq.req_notify()
        qp.post_recv(wr)

    def destroy(self, IDENTITY):
        if IDENTITY not in self.qps:
            return

        mr_ids = [mr_id for mr_id in self.mrs if f"{IDENTITY}-" in mr_id]
        self.qps[IDENTITY].close()
        for mr_id in mr_ids:
            # del self.mrs[mr_id]
            self.mrs.pop(mr_id, None)

        wr_ids = [wr_id for wr_id in self.map if f"{IDENTITY}-" in self.map[wr_id]]
        for wr_id in wr_ids:
            # del self.map[wr_id]
            self.map.pop(wr_id, None)

        if len(self.mrs) == 0:
            self.comp_ch.close()
            self.cq.close()

class DTC:
    def __init__(self, ip, port, nc=100):
        self.IDENTITY = f"{ip}:{port}"
        self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

        self.recv_bk = {}

    def send(self, data_obj, select_socket=None):
        data_b = pickle.dumps(data_obj)

        header = struct.pack("!I", len(data_b))

        send_b = header + data_b

        if select_socket:
            select_socket.sendall(send_b)
        else:
            self.socket.sendall(send_b)

    def recv(self, select_socket=None):
        if not select_socket:
            select_socket = self.socket

        fd = select_socket.fileno()

        if fd in self.recv_bk and self.recv_bk[fd]:
            data_b, len_data, left_len = self.recv_bk[fd]
        else:
            recv_data = select_socket.recv(4)
            if not recv_data:
                return None
            len_data = struct.unpack("!I", recv_data)[0]
            left_len = len_data
            data_b = b''

        while left_len:
            try:
                left_recv = select_socket.recv(left_len)
                recv_len = len(left_recv)
                data_b += left_recv
                left_len -= recv_len
            except BlockingIOError as err:
                break

        if left_len:
            self.recv_bk[fd] = (data_b, len_data, left_len)
            return ''
        else:
            self.recv_bk[fd] = None

        return pickle.loads(data_b)


class Client(DTC):
    def __init__(self, ip, port, nc=100):
        super(Client, self).__init__(ip, port)
        while True:
            return_code = self.socket.connect_ex((ip, port))
            if return_code == 0:
                break
            time.sleep(0.01)

    def clean(self):
        self.socket.close()

class Server(DTC):
    def __init__(self, ip, port, nc=100):
        super(Server, self).__init__(ip, port)
        self.socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        self.socket.bind((ip, port))
        self.socket.listen(nc)

        self.client_socket = []

        self.poller = select.epoll()
        self.map = {}
        self.server_fd = self.socket.fileno()
        self.poller.register(self.server_fd, select.EPOLLIN)
        self.map[self.server_fd] = self.socket
        self.pam = {f"{ip}:{port}": self.server_fd}
        self.address = {}

        self.to_del = []

    def add_dtc(self, dtc):
        dtc_fd = dtc.socket.fileno()
        self.poller.register(dtc_fd, select.EPOLLIN)
        self.map[dtc_fd] = dtc.socket
        self.pam[dtc.IDENTITY] = dtc_fd
        self.address[dtc_fd] = dtc.IDENTITY

    def add_comp_ch(self, comp_ch_fd):
        self.comp_ch = comp_ch_fd
        self.poller.register(comp_ch_fd, select.EPOLLIN)

    def del_dtc(self, dtc_socket):
        dtc_fd = dtc_socket.fileno()
        IDENTITY = [k for k in self.pam if self.pam[k] == dtc_fd][0]
        self.pam.pop(IDENTITY, None)
        self.map.pop(dtc_fd, None)
        self.address.pop(dtc_fd, None)
        # del self.pam[IDENTITY]
        # del self.map[dtc_fd]
        # del self.address[dtc_fd]
        if dtc_socket in self.client_socket:
            self.client_socket.remove(dtc_socket)
        self.poller.unregister(dtc_fd)
        dtc_socket.close()

    def poll_fd(self, fd):
        poll_dtc = self.map[fd]
        recv_data = self.recv(poll_dtc)

        if recv_data is None:
            self.to_del.append(poll_dtc)
            
        return (poll_dtc, recv_data)
        
    def poll(self):
        if self.to_del:
            for dtc_socket in self.to_del:
                self.del_dtc(dtc_socket)
            self.to_del = []

        poll_fds = []
        events = dict(self.poller.poll(timeout=5))
        
        for fd, event in events.items():
            if fd == self.server_fd:
                client_socket, client_address = self.socket.accept()
                client_socket.setblocking(0)
                self.client_socket.append(client_socket)

                self.poller.register(client_socket.fileno(), select.EPOLLIN)
                self.map[client_socket.fileno()] = client_socket
                self.pam[f"{client_address[0]}:{client_address[1]}"] = client_socket.fileno()
                self.address[client_socket.fileno()] = f"{client_address[0]}:{client_address[1]}"
            elif event == select.EPOLLIN:
                poll_fds.append(fd)

        return poll_fds

    def clean(self):
        try:
            for client_socket in self.client_socket:
                if client_socket.fileno() < 0:
                    continue
                self.poller.unregister(client_socket.fileno())
                client_socket.close()
            self.poller.unregister(self.server_fd)
            self.socket.close()
        except Exception as e:
            pass