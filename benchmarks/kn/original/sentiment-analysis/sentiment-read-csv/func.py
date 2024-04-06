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

redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])


def main(context: Context):
    if 'request' in context.keys():
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "0" + "\n", flush=True)
        event = context.request.json

        bucket_name = event['bucket_name']
        file_key = event['file_key']

        response = redisClient.get(file_key)

        lines = response.decode('utf-8').split('\n')

        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "1" + "\n", flush=True)
        for row in csv.DictReader(lines):
            response = requests.get(url=os.environ["SENTIMENT_PRODUCT_OR_SERVICE"], json=row)
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "2" + "\n", flush=True)
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
