import socket
import select
import time
import pickle
import struct

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
        

    def poll(self):
        if self.to_del:
            for dtc_socket in self.to_del:
                self.del_dtc(dtc_socket)
            self.to_del = []
        poll_results = []
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
                # poll_dtcs.append(self.map[fd])
                poll_dtc = self.map[fd]
                recv_data = self.recv(poll_dtc)
                if recv_data:
                    poll_results.append((poll_dtc, recv_data))
                elif recv_data is None:
                    self.to_del.append(poll_dtc)
                    poll_results.append((poll_dtc, None))
                    # self.del_dtc(poll_dtc)

        return poll_results

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