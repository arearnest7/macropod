import os
import select
import zmq
import pickle
import struct

class DTC:
    def __init__(self, read_pipe, write_pipe, ipc_path="/home/app/ipc"):
        self.read_pipe = f"{ipc_path}/{read_pipe}"
        self.write_pipe = f"{ipc_path}/{write_pipe}"
        self.mkfifo(self.read_pipe)
        self.mkfifo(self.write_pipe)

        self.rf = os.open(self.read_pipe, os.O_RDONLY | os.O_NONBLOCK)
        self.wf = os.open(self.write_pipe, os.O_SYNC | os.O_CREAT | os.O_RDWR)

        self.IDENTITY = f"{read_pipe}:{write_pipe}"

    def mkfifo(self, pipe_name):
        try:
            if not os.path.exists(pipe_name):
                os.mkfifo(pipe_name)
        except FileExistsError as err:
            pass

    def send(self, data_obj):
        data_b = pickle.dumps(data_obj)
        header = struct.pack("!I", len(data_b))
        os.write(self.wf, header + data_b)

    def recv(self):
        recv_data_list = []
        while True:
            try:
                header_bytes = b''
                while len(header_bytes) < 4:
                    chunk = os.read(self.rf, 4 - len(header_bytes))
                    if not chunk:
                        return recv_data_list  # EOF
                    header_bytes += chunk

                expected_len = struct.unpack("!I", header_bytes)[0]

                recv_data_b = b''
                while len(recv_data_b) < expected_len:
                    chunk = os.read(self.rf, expected_len - len(recv_data_b))
                    if not chunk:
                        raise EOFError("Early EOF while reading data body")
                    recv_data_b += chunk

                obj = pickle.loads(recv_data_b)
                recv_data_list.append(obj)
            except BlockingIOError as err:
                break
            except Exception as e:
                print(f"[recv error] {e}")
                break

        return recv_data_list

class Client(DTC):
    def __init__(self, pipe_name):
        super(Client, self).__init__(pipe_name)
        self.wf = os.open(self.pipe_name, os.O_SYNC | os.O_CREAT | os.O_RDWR)

class Server(DTC):
    def __init__(self, pipe_name):
        super(Server, self).__init__(pipe_name)
        self.rf = os.open(self.pipe_name, os.O_RDONLY | os.O_NONBLOCK)

class Socket:
    def __init__(self, ip):
        context = zmq.Context()
        self.socket = context.socket(zmq.REQ)
        self.socket.connect(f"tcp://{ip}")

    def send(self, data_obj):
        self.socket.send(pickle.dumps(data_obj))

    def recv(self):
        return pickle.loads(self.socket.recv())
    
class SSocket:
    def __init__(self, ip):
        context = zmq.Context()
        self.socket = context.socket(zmq.REP)
        self.socket.bind(f"tcp://{ip}")

    def send(self, data_obj):
        self.socket.send(pickle.dumps(data_obj))

    def recv(self):
        return pickle.loads(self.socket.recv())

class Poller:
    def __init__(self):
        self.poller = select.epoll()
        self.map = {}

    def add_peer(self, dtc):
        self.poller.register(dtc.rf, select.EPOLLIN | select.EPOLLET)
        self.map[dtc.rf] = dtc

    def del_peer(self, dtc):
        self.poller.unregister(dtc.rf)
        # del self.map[dtc.rf]
        self.map.pop(dtc.rf, None)

    def poll(self):
        events = dict(self.poller.poll(timeout=5))
        dtcs = [self.map[rf] for rf in events]
        return dtcs
    
    def close(self):
        for dtc in self.map.values():
            try:
                os.close(dtc.rf)
                os.close(dtc.wf)
            except Exception as e:
                pass
        self.poller.close()
