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

# purpose: serve the function invoke and return the results to the appropriate request handler
# input:
#  request: MacroPodRequest object containing the payload and all invocation related data
# output:
#  string: reply string
#  int: response code
def Serve_Invoke(request):
    ctx_in = {"Workflow": request.Workflow, "Function": request.Function, "Target": request.Target, "Text": request.Text, "JSON": json_format.MessageToDict(request.JSON), "Data": request.Data, "WorkflowID": request.WorkflowID, "Depth": request.Depth, "Width": request.Width}
    Timestamp(ctx_in, "", "request")
    reply, code = FunctionHandler(ctx_in)
    Timestamp(ctx_in, "", "request_end")
    return reply, code

# purpose: gRPC handler for Invoke
# input:
#  context: context object utilized in every gRPC request handler
#  request: MacroPodRequest object defined by protobuf
# output:
#  MacroPodReply: reply protobuf struct
class MacroPodFunctionServicer(pb_grpc.MacroPodFunctionServicer):
    def Invoke(self, request, context):
        reply, code = Serve_Invoke(request)
        res = pb.MacroPodReply(Reply=reply, Code=code)
        return res

# purpose: http handler for Invoke
# input:
#  res: writer to provide HTTP reply data
#  req: GET/POST object
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

# purpose: deploys the HTTP and gRPC services for function as individual threads
# input:
#  id: int for which service to deploy; both are deployed via main function (got to love weird python jank)
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

# input:
#  SERVICE_PORT: string defining gRPC service binding; default 5000
#  HTTP_PORT: string defining HTTP service binding; default 6000
#  WORKFLOW: string name of workflow the function belongs to
#  FUNCTION: string name of the function
#  INGRESS: string representing the IP:port of the ingress such that gateway communication can be used if enabled
#  LOGGER: string IP:port of the logger if logging is enabled; empty if logging is disabled
#  COMM_TYPE: string name of the communication type (refer to the deployer documentation)
if __name__ == '__main__':
    servers = []
    with futures.ThreadPoolExecutor(max_workers=2) as executor:
        servers.append(executor.submit(server_start, 0))
        servers.append(executor.submit(server_start, 1))
    serving = [server for server in servers]
