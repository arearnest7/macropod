from flask import Flask, request
from flask_restful import Resource, Api
import os
import json
import subprocess
import datetime
from func import function_handler

app = Flask(__name__)
api = Api(app)

class func(Resource):
    def get(self):
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "GET" + "," + "0" + "\n", flush=True)
        if request.is_json:
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "GET" + "," + "1" + "\n", flush=True)
            res = function_handler({"request": request.json, "content-type": "GET", "is_json": True})
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "GET" + "," + "2" + "\n", flush=True)
            return res
        else:
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "GET" + "," + "3" + "\n", flush=True)
            res = function_handler({"request": request.get_data(as_text=True), "content-type": "GET", "is_json": False})
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "GET" + "," + "4" + "\n", flush=True)
            return res
    def post(self):
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "5" + "\n", flush=True)
        if request.is_json:
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "6" + "\n", flush=True)
            res = function_handler({"request": request.json, "content-type": "POST", "is_json": True})
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "7" + "\n", flush=True)
            return res
        else:
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "8" + "\n", flush=True)
            res = function_handler({"request": request.get_data(as_text=True), "content-type": "POST", "is_json": False})
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "9" + "\n", flush=True)
            return res

api.add_resource(func, '/')

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=int(os.environ["PORT"]))
