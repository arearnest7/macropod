from rpc import RPC
import base64
import json
import os
import sys
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def function_handler(context):
    if context["InvokeType"] == "GRPC":
        params = json.loads(context["Request"])

        #redisClient.set("merit-" + str(params["id"]), json.dumps(params))

        return str(params["id"]) + " statistics uploaded/updated", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
