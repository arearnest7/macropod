from parliament import Context
from flask import Request
import base64
import requests
import csv
import json
import os
import datetime
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])


def main(context: Context):
    if 'request' in context.keys():
        event = context.request.json

        bucket_name = event['bucket_name']
        file_key = event['file_key']

        response = open(file_key, 'r').read()

        lines = response.split('\n')

        for row in csv.DictReader(lines):
            response = requests.post(url=os.environ["SENTIMENT_PRODUCT_OR_SERVICE"], json=row)
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
