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


def main(context: Context):
    if 'request' in context.keys():
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "0" + "\n", flush=True)
        event = context.request.json["event"]
        data = redisClient.get("original-" + event[0])
        md5 = hashlib.md5(data).hexdigest()
        if event[1] == md5:
            redisClient.set("checksumed-" + event[0], data)
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "1" + "\n", flush=True)
            return "success", 200
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "2" + "\n", flush=True)
        return "failed", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
