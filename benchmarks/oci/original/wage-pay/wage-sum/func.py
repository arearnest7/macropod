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
        temp = json.loads(open(params["operator"], 'r').read())
        params["operator"] = temp["operator"]
        params["id"] = temp["id"]
        stats = {'total': params['total']['statistics']['total'] }
        params['statistics'] = stats

        response = requests.get(url=os.environ["WAGE_AVG"], json=params)
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
