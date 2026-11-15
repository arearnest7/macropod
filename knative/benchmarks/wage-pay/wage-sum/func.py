from parliament import Context
from flask import Request
import base64
import requests
import json
import os
import sys
import datetime
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])


def main(context: Context):
    if 'request' in context.keys():
        params = context.request.json
        temp = json.loads(open(params["operator"], 'r').read())
        params["operator"] = temp["operator"]
        params["id"] = temp["id"]
        stats = {'total': params['total']['statistics']['total'] }
        params['statistics'] = stats

        response = requests.post(url=os.environ["WAGE_AVG"], json=params)
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
