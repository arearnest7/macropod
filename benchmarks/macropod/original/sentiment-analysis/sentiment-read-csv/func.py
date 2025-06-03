from rpc import Invoke_JSON
import base64
import requests
import csv
import json
import os
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def FunctionHandler(context):
    event = json.loads(context["Text"])

    bucket_name = event['bucket_name']
    file_key = event['file_key']

    response = open(file_key, 'r').read()

    lines = response.split('\n')

    for row in csv.DictReader(lines):
        response = Invoke_JSON(context, "SENTIMENT_PRODUCT_OR_SERVICE", row])[0]
    return response, 200
