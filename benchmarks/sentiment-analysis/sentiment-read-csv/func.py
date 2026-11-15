from rpc import Invoke_JSON
import base64
import csv
import json
import os
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def FunctionHandler(context):
    event = context["JSON"]

    bucket_name = event['bucket_name']
    file_key = event['file_key']

    response = open(file_key, 'r').read()

    lines = response.split('\n')

    results = []

    for row in csv.DictReader(lines):
        results = Invoke_JSON(context, "SENTIMENT_PRODUCT_OR_SERVICE", row)
    return results[0], 200
