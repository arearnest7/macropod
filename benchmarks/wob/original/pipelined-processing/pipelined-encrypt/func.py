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

if "LOGGING_NAME" in os.environ:
    loggingClient = redis.Redis(host=os.environ['LOGGING_URL'], password=os.environ['LOGGING_PASSWORD'])

def main(context: Context):
    if 'request' in context.keys():
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n")
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
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n")
        return "success", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
