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
import app_pb2 as pb
import app_pb2_grpc as pb_grpc
from func import function_handler

app = Flask(__name__)
api = Api(app)

if "LOGGING_NAME" in os.environ:
    redisClient = redis.Redis(host=os.environ['LOGGING_URL'], password=os.environ['LOGGING_PASSWORD'])

class HTTPFunctionHandler(Resource):
    def get(self):
        if request.is_json:
            workflow_id = str(random.randint(0, 10000000))
            if "LOGGING_NAME" in os.environ:
                redisClient.append(os.environ["LOGGING_NAME"], "EXECUTION_GET_JSON_START," + workflow_id + ",NA,1," + str(datetime.datetime.now()) + "\n")
            reply, code = function_handler({"request": request.json, "workflow_id": workflow_id, "request_type": "GET", "is_json": True})
            if "LOGGING_NAME" in os.environ:
                redisClient.append(os.environ["LOGGING_NAME"], "EXECUTION_GET_JSON_END," + workflow_id + ",NA,1," + str(datetime.datetime.now()) + "\n")
            return reply
        else:
            workflow_id = str(random.randint(0, 10000000))
            if "LOGGING_NAME" in os.environ:
                redisClient.append(os.environ["LOGGING_NAME"], "EXECUTION_GET_TEXT_START," + workflow_id + ",NA,1," + str(datetime.datetime.now()) + "\n")
            reply, code = function_handler({"request": {"text": request.get_data(as_text=True)}, "workflow_id": workflow_id, "request_type": "GET", "is_json": False})
            if "LOGGING_NAME" in os.environ:
                redisClient.append(os.environ["LOGGING_NAME"], "EXECUTION_GET_TEXT_END," + workflow_id + ",NA,1," + str(datetime.datetime.now()) + "\n")
            return reply
    def post(self):
        if request.is_json:
            workflow_id = str(random.randint(0, 10000000))
            if "LOGGING_NAME" in os.environ:
                redisClient.append(os.environ["LOGGING_NAME"], "EXECUTION_POST_JSON_START," + workflow_id + ",NA,1," + str(datetime.datetime.now()) + "\n")
            reply, code = function_handler({"request": request.json, "workflow_id": workflow_id, "request_type": "POST", "is_json": True})
            if "LOGGING_NAME" in os.environ:
                redisClient.append(os.environ["LOGGING_NAME"], "EXECUTION_POST_JSON_END," + workflow_id + ",NA,1," + str(datetime.datetime.now()) + "\n")
            return reply
        else:
            workflow_id = str(random.randint(0, 10000000))
            if "LOGGING_NAME" in os.environ:
                redisClient.append(os.environ["LOGGING_NAME"], "EXECUTION_POST_TEXT_START," + workflow_id + ",NA,1," + str(datetime.datetime.now()) + "\n")
            reply, code = function_handler({"request": {"text": request.get_data(as_text=True)}, "workflow_id": workflow_id, "request_type": "POST", "is_json": False})
            if "LOGGING_NAME" in os.environ:
                redisClient.append(os.environ["LOGGING_NAME"], "EXECUTION_POST_TEXT_END," + workflow_id + ",NA,1," + str(datetime.datetime.now()) + "\n")
            return reply

api.add_resource(HTTPFunctionHandler, '/')

class gRPCFunctionServicer(pb_grpc.gRPCFunctionServicer):
    def gRPCFunctionHandler(self, request, context):
        if "LOGGING_NAME" in os.environ:
            redisClient.append(os.environ["LOGGING_NAME"], "EXECUTION_GRPC_START," + request.workflow_id + ",NA,1," + str(datetime.datetime.now()) + "\n")
        reply, code = function_handler({"request": request.body, "workflow_id": request.workflow_id, "request_type": "GRPC", "is_json": False})
        res = pb.ResponseBody(reply=reply, code=code)
        if "LOGGING_NAME" in os.environ:
            redisClient.append(os.environ["LOGGING_NAME"], "EXECUTION_GRPC_END," + request.workflow_id + ",NA,1," + str(datetime.datetime.now()) + "\n")
        return res

if __name__ == '__main__':
    if "SERVICE_TYPE" not in os.environ or os.environ["SERVICE_TYPE"] == "HTTP":
        app.run(host='0.0.0.0', port=int(os.environ["FUNC_PORT"]))
    elif os.environ["SERVICE_TYPE"] == "GRPC":
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=int(os.environ["GRPC_THREAD"])))
        pb_grpc.add_gRPCFunctionServicer_to_server(gRPCFunctionServicer(), server)
        server.add_insecure_port("[::]:" + os.environ["FUNC_PORT"])
        server.start()
        server.wait_for_termination()
