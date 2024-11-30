from rpc import RPC
import base64
import requests
import csv
import json
import os
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def FunctionHandler(context):
    event = json.loads(context["Request"])

    bucket_name = event['bucket_name']
    file_key = event['file_key']

    response = open(file_key, 'r').read()

    lines = response.split('\n')

    for row in csv.DictReader(lines):
        payload = []
        payload.append(json.dumps(row).encode())
        response = RPC(context, os.environ["SENTIMENT_PRODUCT_OR_SERVICE"], payload)[0]
    return response, 200
