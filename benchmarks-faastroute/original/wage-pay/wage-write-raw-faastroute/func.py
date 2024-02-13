import base64
import requests
import json
import os
import sys
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def function_handler(context):
    if context["is_json"]:
        params = context["request"]
        #redisClient.set("raw-" + str(params["id"]), json.dumps(params))
        response = requests.get(url=os.environ["WAGE_STATS"], json={})
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
