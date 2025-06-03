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
        #redisClient.set("raw-" + str(params["id"]), json.dumps(params))
        response = requests.post(url=os.environ["WAGE_STATS"], json={})
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
