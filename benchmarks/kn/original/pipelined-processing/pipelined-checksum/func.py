from parliament import Context
from flask import Request
import base64
import os
import json
import string
import hashlib
import datetime
import redis
import random

redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

if "LOGGING_NAME" in os.environ:
    loggingClient = redis.Redis(host=os.environ['LOGGING_URL'], password=os.environ['LOGGING_PASSWORD'])

def main(context: Context):
    if 'request' in context.keys():
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n")
        event = context.request.json["event"]
        data = redisClient.get("original-" + event[0])
        md5 = hashlib.md5(data).hexdigest()
        if event[1] == md5:
            redisClient.set("checksumed-" + event[0], data)
            if "LOGGING_NAME" in os.environ:
                loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n")
            return "success", 200
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "2" + "\n")
        return "failed", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
