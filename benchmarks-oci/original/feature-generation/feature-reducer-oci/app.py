from flask import Flask, request
from flask_restful import Resource, Api
import os
import json
import subprocess
from func import function_handler

app = Flask(__name__)
api = Api(app)

class func(Resource):
    def get(self):
        if request.is_json:
            return function_handler({"request": request.json, "content-type": "GET", "is_json": True})
        else:
            return function_handler({"request": request.get_data(as_text=True), "content-type": "GET", "is_json": False})
    def post(self):
        if request.is_json:
            return function_handler({"request": request.json, "content-type": "POST", "is_json": True})
        else:
            return function_handler({"request": request.get_data(as_text=True), "content-type": "POST", "is_json": False})

api.add_resource(func, '/')

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=int(os.environ["PORT"]))
