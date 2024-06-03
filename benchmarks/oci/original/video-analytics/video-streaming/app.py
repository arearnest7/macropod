from flask import Flask, request
from flask_restful import Resource, Api
import os
import json
import subprocess
import datetime
import random
from func import function_handler

app = Flask(__name__)
api = Api(app)

class func(Resource):
    def get(self):
        workflow_id = str(random.randint(0, 10000000))
        workflow_depth = 0
        workflow_width = 0
        body = request.json
        if "workflow_id" in body:
            workflow_id = body["workflow_id"]
            workflow_depth = body["workflow_depth"]
            workflow_width = body["workflow_width"]
        else:
            body["workflow_id"] = workflow_id
            body["workflow_depth"] = workflow_depth
            body["workflow_width"] = workflow_width
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "0" + "\n", flush=True)
        if request.is_json:
            print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "1" + "\n", flush=True)
            res = function_handler({"request": body, "content-type": "HTTP", "is_json": True})
            print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "2" + "\n", flush=True)
            return res
        else:
            print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "1" + "\n", flush=True)
            res = function_handler({"request": body, "content-type": "HTTP", "is_json": False})
            print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "2" + "\n", flush=True)
            return res
    def post(self):
        workflow_id = str(random.randint(0, 10000000))
        workflow_depth = 0
        workflow_width = 0
        body = request.json
        if "workflow_id" in body:
            workflow_id = body["workflow_id"]
            workflow_depth = body["workflow_depth"]
            workflow_width = body["workflow_width"]
        else:
            body["workflow_id"] = workflow_id
            body["workflow_depth"] = workflow_depth
            body["workflow_width"] = workflow_width
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "0" + "\n", flush=True)
        if request.is_json:
            print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "1" + "\n", flush=True)
            res = function_handler({"request": body, "content-type": "HTTP", "is_json": True})
            print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "2" + "\n", flush=True)
            return res
        else:
            print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "1" + "\n", flush=True)
            res = function_handler({"request": body, "content-type": "HTTP", "is_json": False})
            print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "2" + "\n", flush=True)
            return res

api.add_resource(func, '/')

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=int(os.environ["PORT"]))
