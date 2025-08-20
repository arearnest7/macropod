from rpc import Invoke
import base64
from zipfile import ZipFile
import os
import json
import string
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def FunctionHandler(context):
    event = context["Text"]
    data = open("checksumed-" + event, 'rb').read()
    with open("/tmp/" + event, "wb") as f:
        f.write(data)
    with ZipFile('/tmp/zip.zip', 'w') as zip:
        zip.write("/tmp/" + event)
    zip.close()
    with open("/tmp/zip.zip", "rb") as f:
        data = f.read()
    #redisClient.set("ziped-" + event, data)
    return "success", 200
