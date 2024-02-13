import base64
import os
import json
import string
import hashlib
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def function_handler(context):
    if context["is_json"]:
        event = context["request"]["event"]
        data = open("original-" + event[0], 'rb').read()
        md5 = hashlib.md5(data).hexdigest()
        if event[1] == md5:
            #redisClient.set("checksumed-" + event[0], data)
            return "success", 200
        return "failed", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
