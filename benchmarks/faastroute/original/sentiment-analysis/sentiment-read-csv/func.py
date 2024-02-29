from rpc import RPC
import base64
import requests
import csv
import json
import os
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def function_handler(context):
    if context["request_type"] == "GRPC":
        event = json.loads(context["request"])

        bucket_name = event['bucket_name']
        file_key = event['file_key']

        response = open(file_key, 'r').read()

        lines = response.split('\n')

        for row in csv.DictReader(lines):
            response = RPC(context, os.environ["SENTIMENT_PRODUCT_OR_SERVICE"], [json.dumps(row)])[0]
        return response, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
