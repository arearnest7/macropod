from flask import Flask, request
from flask_restful import Resource, Api
from concurrent import futures
import os
import json
import random
import subprocess
import datetime
import redis
import grpc
import mmap
import app_pb2 as pb
import app_pb2_grpc as pb_grpc
from func import FunctionHandler

app = Flask(__name__)
api = Api(app)

MAX_MESSAGE_LENGTH = 1024 * 1024 * 200
opts = [("grpc.max_receive_message_length", MAX_MESSAGE_LENGTH),("grpc.max_send_message_length", MAX_MESSAGE_LENGTH)]
if "LOGGING_NAME" in os.environ:
    redisClient = redis.Redis(host=os.environ['LOGGING_URL'], password=os.environ['LOGGING_PASSWORD'])

class HTTPFunctionHandler(Resource):
    def get(self):
        request_type = ["gg", "gm"]["APP_PV" in os.environ]
        if request.is_json:
            workflow_id = str(random.randint(0, 10000000))
            if "LOGGING_NAME" in os.environ:
                redisClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + workflow_id + "," + "0" + "," + "0" + "," + request_type + "," + "0" + "\n")
            reply, code = FunctionHandler({"Request": request.json, "WorkflowId": workflow_id, "Depth": 0, "Width": 0, "RequestType": request_type, "InvokeType": "GET", "IsJson": True})
            if "LOGGING_NAME" in os.environ:
                redisClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + workflow_id + "," + "0" + "," + "0" + "," + request_type + "," + "1" + "\n")
            return reply
        else:
            workflow_id = str(random.randint(0, 10000000))
            if "LOGGING_NAME" in os.environ:
                redisClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + workflow_id + "," + "0" + "," + "0" + "," + request_type + "," + "2" + "\n")
            reply, code = FunctionHandler({"Request": request.get_data(as_text=True), "WorkflowId": workflow_id, "Depth": 0, "Width": 0, "RequestType": request_type, "InvokeType": "GET", "IsJson": False})
            if "LOGGING_NAME" in os.environ:
                redisClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + workflow_id + "," + "0" + "," + "0" + "," + request_type + "," + "3" + "\n")
            return reply
    def post(self):
        request_type = ["gg", "gm"]["APP_PV" in os.environ]
        if request.is_json:
            workflow_id = str(random.randint(0, 10000000))
            if "LOGGING_NAME" in os.environ:
                redisClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + workflow_id + "," + "0" + "," + "0" + "," + request_type + "," + "4" + "\n")
            reply, code = FunctionHandler({"Request": request.json, "WorkflowId": workflow_id, "Depth": 0, "Width": 0, "RequestType": request_type, "InvokeType": "POST", "IsJson": True})
            if "LOGGING_NAME" in os.environ:
                redisClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + workflow_id + "," + "0" + "," + "0" + "," + request_type + "," + "5" + "\n")
            return reply
        else:
            workflow_id = str(random.randint(0, 10000000))
            if "LOGGING_NAME" in os.environ:
                redisClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + workflow_id + "," + "0" + "," + "0" + "," + request_type + "," + "6" + "\n")
            reply, code = FunctionHandler({"Request": request.get_data(as_text=True), "WorkflowId": workflow_id, "Depth": 0, "Width": 0, "RequestType": request_type, "InvokeType": "POST", "IsJson": False})
            if "LOGGING_NAME" in os.environ:
                redisClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + workflow_id + "," + "0" + "," + "0" + "," + request_type + "," + "7" + "\n")
            return reply

api.add_resource(HTTPFunctionHandler, '/')

class gRPCFunctionServicer(pb_grpc.gRPCFunctionServicer):
    def gRPCFunctionHandler(self, request, context):
        if "LOGGING_NAME" in os.environ:
            redisClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + request.workflow_id + "," + str(request.depth) + "," + str(request.width) + "," + request.request_type + "," + "8" + "\n")
        if request.request_type == "" or request.request_type == "gg":
            print("gg")
            reply, code = FunctionHandler({"Request": request.data, "WorkflowId": request.workflow_id, "Depth": request.depth, "Width": request.width, "RequestType": request.request_type, "InvokeType": "GRPC", "IsJson": False})
            res = pb.ResponseBody(reply=reply, code=code)
        elif request.request_type == "mg":
            print("mg")
            req = b''
            with open(os.environ["APP_PV"] + "/" + request.pv_path, "rb") as f:
                mm = mmap.mmap(f.fileno(), 0, access=mmap.ACCESS_READ)
                req = mm.read()
                mm.close()
            reply, code = FunctionHandler({"Request": req, "WorkflowId": request.workflow_id, "Depth": request.depth, "Width": request.width, "RequestType": request.request_type, "InvokeType": "GRPC", "IsJson": False})
            res = pb.ResponseBody(reply=reply, code=code)
        elif request.request_type == "gm":
            print("gm")
            payload, code = FunctionHandler({"Request": request.data, "WorkflowId": request.workflow_id, "Depth": request.depth, "Width": request.width, "RequestType": request.request_type, "InvokeType": "GRPC", "IsJson": False})
            pv_path = request.workflow_id + "_" + str(request.depth) + "_" + str(request.width) + "_" + str(random.randint(0, 10000000))
            with open(os.environ["APP_PV"] + "/" + pv_path, "wb") as f:
                f.write(payload)
            reply = ""
            res = pb.ResponseBody(reply=reply, code=code, pv_path=pv_path)
        elif request.request_type == "mm":
            print("mm")
            req = b''
            with open(os.environ["APP_PV"] + "/" + request.pv_path, "rb") as f:
                mm = mmap.mmap(f.fileno(), 0, access=mmap.ACCESS_READ)
                req = mm.read()
                mm.close()
            payload, code = FunctionHandler({"Request": req, "WorkflowId": request.workflow_id, "Depth": request.depth, "Width": request.width, "RequestType": request.request_type, "InvokeType": "GRPC", "IsJson": False})
            pv_path = request.workflow_id + "_" + str(request.depth) + "_" + str(request.width) + "_" + str(random.randint(0, 10000000))
            print(pv_path)
            with open(os.environ["APP_PV"] + "/" + pv_path, "wb") as f:
                f.write(payload)
            reply = ""
            print("reply: " + reply)
            print("code: " + str(code))
            print("pv_path: " + pv_path)
            res = pb.ResponseBody(reply=reply, code=code, pv_path=pv_path)
        else:
            res = pb.ResponseBody(reply="", code=500)
        if "LOGGING_NAME" in os.environ:
            redisClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + request.workflow_id + "," + str(request.depth) + "," + str(request.width) + "," + request.request_type + "," + "9" + "\n")
        return res

if __name__ == '__main__':
    if "SERVICE_TYPE" not in os.environ or os.environ["SERVICE_TYPE"] == "HTTP":
        app.run(host='0.0.0.0', port=int(os.environ["FUNC_PORT"]))
    elif os.environ["SERVICE_TYPE"] == "GRPC":
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=int(os.environ["GRPC_THREAD"])), options=opts)
        pb_grpc.add_gRPCFunctionServicer_to_server(gRPCFunctionServicer(), server)
        server.add_insecure_port("[::]:" + os.environ["FUNC_PORT"])
        server.start()
        server.wait_for_termination()
