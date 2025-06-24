from concurrent import futures
import os
import datetime
import random
import grpc
from flask import Flask, request
from flask_restful import Resource, Api
import macropod_pb2 as pb
import macropod_pb2_grpc as pb_grpc
from func import FunctionHandler
from rpc import Timestamp
from google.protobuf.struct_pb2 import Struct
from google.protobuf import json_format

app = Flask(__name__)
api = Api(app)

MAX_MESSAGE_LENGTH = 1024 * 1024 * 200
opts = [("grpc.max_receive_message_length", MAX_MESSAGE_LENGTH),("grpc.max_send_message_length", MAX_MESSAGE_LENGTH)]

def Serve_Invoke(request):
    ctx_in = {"Workflow": request.Workflow, "Function": request.Function, "Target": request.Target, "Text": request.Text, "JSON": json_format.MessageToDict(request.JSON), "Data": request.Data, "WorkflowID": request.WorkflowID, "Depth": request.Depth, "Width": request.Width}
    Timestamp(ctx_in, "", "request")
    reply, code = FunctionHandler(ctx_in)
    Timestamp(ctx_in, "", "request_end")
    return reply, code

class MacroPodFunctionServicer(pb_grpc.MacroPodFunctionServicer):
    def Invoke(self, request, context):
        reply, code = Serve_Invoke(request)
        res = pb.MacroPodReply(Reply=reply, Code=code)
        return res

class HTTP_Invoke(Resource):
    def get(self):
        workflow = ""
        if "WORKFLOW" in os.environ:
            workflow = os.environ["WORKFLOW"]
        function = ""
        if "FUNCTION" in os.environ:
            function = os.environ["FUNCTION"]
        workflow_id = str(random.randint(0, 100000))
        dw = 0
        target = request.url
        ctx_in = pb.MacroPodRequest(Workflow=workflow, Function=function, WorkflowID=workflow_id, Depth=dw, Width=dw, Target=target)
        content = request.content_type
        if content == "text/plain":
            t = request.data.decode()
            ctx_in.Text = t
        elif content == "application/json":
            j = request.get_json()
            ctx_in.JSON = j
        elif content == "application/octet-stream":
            d = request.data
            ctx_in.Data = d
        else:
            t = request.data.decode()
            ctx_in.Text = t
        reply, code = Serve_Invoke(ctx_in)
        return reply
    def post(self):
        workflow = ""
        if "WORKFLOW" in os.environ:
            workflow = os.environ["WORKFLOW"]
        function = ""
        if "FUNCTION" in os.environ:
            function = os.environ["FUNCTION"]
        workflow_id = str(random.randint(0, 100000))
        dw = 0
        target = request.url
        ctx_in = pb.MacroPodRequest(Workflow=workflow, Function=function, WorkflowID=workflow_id, Depth=dw, Width=dw, Target=target)
        content = request.content_type
        if content == "text/plain":
            t = request.data.decode()
            ctx_in.Text = t
        elif content == "application/json":
            j = request.get_json()
            ctx_in.JSON = j
        elif content == "application/octet-stream":
            d = request.data
            ctx_in.Data = d
        else:
            t = request.data.decode()
            ctx_in.Text = t
        reply, code = Serve_Invoke(ctx_in)
        return reply

api.add_resource(HTTP_Invoke, '/')

def server_start(id):
    if id == 0:
        service_port = os.environ["SERVICE_PORT"]
        if service_port == "":
            service_port = "5000"
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=100), options=opts)
        pb_grpc.add_MacroPodFunctionServicer_to_server(MacroPodFunctionServicer(), server)
        server.add_insecure_port("[::]:" + service_port)
        server.start()
        server.wait_for_termination()
    elif id == 1:
        http_port = os.environ["HTTP_PORT"]
        if http_port == "":
            http_port = "6000"
        app.run(host='0.0.0.0', port=int(http_port))

if __name__ == '__main__':
    servers = []
    with futures.ThreadPoolExecutor(max_workers=2) as executor:
        servers.append(executor.submit(server_start, 0))
        servers.append(executor.submit(server_start, 1))
    serving = [server for server in servers]
