import base64
import os
import json
import string
from cryptography.fernet import Fernet
import redis
import random
import datetime

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def function_handler(context):
    if context["is_json"]:
        workflow_id = context["request"]["workflow_id"]
        workflow_depth = context["request"]["workflow_depth"]
        workflow_width = context["request"]["workflow_width"]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "3" + "\n", flush=True)
        event = context["request"]["event"]
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
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "4" + "\n", flush=True)
        return "success", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
