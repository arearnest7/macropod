from rpc import RPC
import base64
import requests
import json
import os
import sys
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def function_handler(context):
    if context["InvokeType"] == "GRPC":
        params = json.loads(context["Request"])
        #redisClient.set("raw-" + str(params["id"]), json.dumps(params))
        response = RPC(context, os.environ["WAGE_STATS"], ["{}".encode()])[0]
        return response, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
