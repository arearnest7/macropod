from parliament import Context
from flask import Request
import base64
import json
import os
import sys
import redis
import random

redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def main(context: Context):
    if 'request' in context.keys():
        params = context.request.json

        redisClient.set("merit-" + str(params["id"]), json.dumps(params))

        return str(params["id"]) + " statistics uploaded/updated", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
