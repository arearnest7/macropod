from rpc import RPC
import base64
from zipfile import ZipFile
import os
import json
import string
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def FunctionHandler(context):
    event = json.loads(context["Request"])
    data = open("checksumed-" + event[0], 'rb').read()
    with open("/tmp/" + event[0], "wb") as f:
        f.write(data)
    with ZipFile('/tmp/zip.zip', 'w') as zip:
        zip.write("/tmp/" + event[0])
    zip.close()
    with open("/tmp/zip.zip", "rb") as f:
        data = f.read()
    #redisClient.set("ziped-" + event[0], data)
    return "success", 200
