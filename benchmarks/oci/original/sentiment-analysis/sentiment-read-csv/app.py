from flask import Flask, request
from flask_restful import Resource, Api
import os
import json
import subprocess
import datetime
import redis
from func import function_handler

app = Flask(__name__)
api = Api(app)

MAX_MESSAGE_LENGTH = 1024 * 1024 * 200
opts = [("grpc.max_receive_message_length", MAX_MESSAGE_LENGTH),("grpc.max_send_message_length", MAX_MESSAGE_LENGTH)]
if "LOGGING_NAME" in os.environ:
    redisClient = redis.Redis(host=os.environ['LOGGING_URL'], password=os.environ['LOGGING_PASSWORD'])

class func(Resource):
    def get(self):
        if "LOGGING_NAME" in os.environ:
            redisClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "http" + "," + "0" + "\n")
        if request.is_json:
            return function_handler({"request": request.json, "content-type": "GET", "is_json": True})
        else:
            return function_handler({"request": request.get_data(as_text=True), "content-type": "GET", "is_json": False})
        if "LOGGING_NAME" in os.environ:
            redisClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "http" + "," + "1" + "\n")
    def post(self):
        if "LOGGING_NAME" in os.environ:
            redisClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "http" + "," + "2" + "\n")
        if request.is_json:
            return function_handler({"request": request.json, "content-type": "POST", "is_json": True})
        else:
            return function_handler({"request": request.get_data(as_text=True), "content-type": "POST", "is_json": False})
        if "LOGGING_NAME" in os.environ:
            redisClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "http" + "," + "3" + "\n")

api.add_resource(func, '/')

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=int(os.environ["PORT"]))
