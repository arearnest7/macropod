from parliament import Context
from flask import Request
import base64
import os
import json
import string
from cryptography.fernet import Fernet
import datetime
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])


def main(context: Context):
    if 'request' in context.keys():
        workflow_id = str(random.randint(0, 10000000))
        workflow_depth = 0
        workflow_width = 0
        if "workflow_id" in context.request.json:
            workflow_id = context.request.json["workflow_id"]
            workflow_depth = context.request.json["workflow_depth"]
            workflow_width = context.request.json["workflow_width"]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "0" + "\n", flush=True)
        event = context.request.json["event"]
        data = open("ziped-" + event[0], 'rb').read()
        with open("/tmp/" + event[0] + ".zip", "wb") as f:
            f.write(data)
        key = Fernet.generate_key()
        with open('/tmp/key.key', 'wb') as filekey:
            filekey.write(key)
        filekey.close()
        fernet = Fernet(key)
        data = ""
        with open("/tmp/" + event[0] + ".zip", "rb") as file:
            data = file.read()
        file.close()
        encrypted_data = fernet.encrypt(data)
        #redisClient.set("encrypted-" + event[0], encrypted_data)
        #redisClient.set("encrypted-key-" + event[0], key)
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "1" + "\n", flush=True)
        return "success", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
