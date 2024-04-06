import base64
import os
import json
import string
import hashlib
import redis
import random
import datetime

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def function_handler(context):
    if context["is_json"]:
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "10" + "\n", flush=True)
        event = context["request"]["event"]
        data = open("original-" + event[0], 'rb').read()
        md5 = hashlib.md5(data).hexdigest()
        if event[1] == md5:
            #redisClient.set("checksumed-" + event[0], data)
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "11" + "\n", flush=True)
            return "success", 200
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "12" + "\n", flush=True)
        return "failed", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
