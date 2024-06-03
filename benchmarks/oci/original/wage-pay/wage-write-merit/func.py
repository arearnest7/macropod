import base64
import json
import os
import sys
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
        params = context["request"]

        #redisClient.set("merit-" + str(params["id"]), json.dumps(params))

        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "4" + "\n", flush=True)
        return str(params["id"]) + " statistics uploaded/updated", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
