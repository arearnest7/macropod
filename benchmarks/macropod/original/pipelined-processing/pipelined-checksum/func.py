from rpc import Invoke
import base64
import os
import json
import string
import hashlib
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def FunctionHandler(context):
    event = context["JSON"]
    data = open("original-" + event[0], 'rb').read()
    md5 = hashlib.md5(data).hexdigest()
    if event[1] == md5:
        #redisClient.set("checksumed-" + event[0], data)
        return "success", 200
    return "failed", 200
