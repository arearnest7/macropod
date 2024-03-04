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
        temp = json.loads(open(params["operator"], 'r').read())
        params["operator"] = temp["operator"]
        params["id"] = temp["id"]
        stats = {'total': params['total']['statistics']['total'] }
        params['statistics'] = stats

        response = RPC(context, os.environ["WAGE_AVG"], [json.dumps(params).encode()])[0]
        return response, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
