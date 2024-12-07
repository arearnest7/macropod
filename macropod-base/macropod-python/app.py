from flask import Flask, request
from flask_restful import Resource, Api
from concurrent import futures
import os
import datetime
import json
from func import FunctionHandler

app = Flask(__name__)
api = Api(app)

class function(Resource):
    def post(self):
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + request.json["WorkflowId"] + "," + str(request.json["Depth"]) + "," + str(request.json["Width"]) + "," + request.json["RequestType"] + "," + "request_start" + "\n", flush=True)
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + request.json["WorkflowId"] + "," + str(request.json["Depth"]) + "," + str(request.json["Width"]) + "," + request.json["RequestType"] + "," + "function_start" + "\n", flush=True)
        reply, code = FunctionHandler({"Request": request.json["Data"].encode(), "WorkflowId": request.json["WorkflowId"], "Depth": request.json["Depth"], "Width": request.json["Width"], "RequestType": request.json["RequestType"], "InvokeType": "GRPC", "IsJson": False})
        res = reply
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + request.json["WorkflowId"] + "," + str(request.json["Depth"]) + "," + str(request.json["Width"]) + "," + request.json["RequestType"] + "," + "request_end" + "\n", flush=True)
        return res

api.add_resource(function, '/')

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=int(os.environ["FUNC_PORT"]))
