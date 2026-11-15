import json
import uuid
import time

from dtc import Client
from message_types import *

def generate_data():
    # func for generate req for message
    return ""

def generate_reqs(n_req):
    reqs = []
    for i in range(n_req):
        data = generate_data()
        req_id = str(uuid.uuid4())
        message = Message(type_id=MESSAGE_DATA, req_id=req_id, req=data, index=0)
        reqs.append(message)
    return reqs

def get_socket(cc_name):
    cc_ip = f"{cc_name}.openfaas-fn.svc.cluster.local"
    cc_port = 6000
    socket = Client(cc_ip, cc_port)
    return socket

def get_entry(socket, func_name):
    message = Message(type_id=MESSAGE_ENTRY, req_id=None, req=func_name, index=None)
    socket.send(message)
    response = socket.recv()

    clients = []
    for client_index, client_ip in response.req:
        clients.append(Client(client_ip, 6000))

    return clients

def send_reqs(clients, st, reqs):
    for i, req in enumerate(reqs):
        client = clients[i % len(clients)]
        client.send(req)
        time.sleep(st)

def get_res(socket):
    message = Message(type_id=MESSAGE_LAT_DIS, req_id=None, req=None, index=None)
    socket.send(message)
    response = socket.recv()

    return response.req

def clear(socket):
    message = Message(type_id=MESSAGE_CLEAR, req_id=None, req=None, index=None)
    socket.send(message)
    socket.recv()

def close_dtc(clients):
    for dtc in clients:
        try:
            dtc.socket.close()
        except Exception as e:
            pass

def handle(req):
    """handle a request to the function
    Args:
        req (str): request body
    """
    try:
        req_j = json.loads(req)
        n_req = req_j["n"]
        st = req_j["st"]
    except Exception as e:
        n_req = 11
        st = 1

    reqs = generate_reqs(n_req)

    cc_name = "your_cc_name"

    socket = get_socket(cc_name)

    entry_func = "your_entry_func"

    clients = get_entry(socket, entry_func)

    send_reqs(clients, st, reqs)

    res = get_res(socket)

    clear(socket)

    socket.socket.close()

    close_dtc(clients)

    return res
