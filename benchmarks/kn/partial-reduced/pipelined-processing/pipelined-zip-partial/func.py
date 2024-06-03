from parliament import Context
from flask import Request
import base64
from zipfile import ZipFile
import os
import json
import string
import datetime
import redis
import random

redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])


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
        data = redisClient.get("checksumed-" + event[0])
        with open("/tmp/" + event[0], "wb") as f:
            f.write(data)
        with ZipFile('/tmp/zip.zip', 'w') as zip:
            zip.write("/tmp/" + event[0])
        zip.close()
        with open("/tmp/zip.zip", "rb") as f:
            data = f.read()
        redisClient.set("ziped-" + event[0], data)
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "1" + "\n", flush=True)
        return "success", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
