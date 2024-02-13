from parliament import Context
from flask import Request
import base64
import requests
import json
import pprint
import csv
import os
import redis
import random

pp = pprint.PrettyPrinter(indent=4)
redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def cfail_handler(req):
    return "CategorizationFail: Fail: \"Input CSV could not be categorised into 'Product' or 'Service'.\""

def product_or_service_handler(req):
    event = req

    results = ""
    if event["reviewType"] == "Product":
        response = requests.get(url=os.environ["SENTIMENT_PRODUCT_SENTIMENT_PRS"], json=event)
        results = response.text
    elif event["reviewType"] == "Service":
        response = requests.get(url=os.environ["SENTIMENT_SERVICE_SENTIMENT_SRS"], json=event)
        results = response.text
    else:
        results = cfail_handler(event)
    return results

def read_csv_handler(req):
    event = req

    bucket_name = event['bucket_name']
    file_key = event['file_key']

    response = redisClient.get(file_key)

    lines = response.decode('utf-8').split('\n')

    for row in csv.DictReader(lines):
        return product_or_service_handler(row)

def main(context: Context):
    if 'request' in context.keys():
        event = context.request.json

        try:
            pp
        except NameError:
            pp = pprint.PrettyPrinter(indent=4)

        bucket_name = event['Records'][0]['s3']['bucket']['name']
        file_key = event['Records'][0]['s3']['object']['key']


        input= {
                'bucket_name': bucket_name,
                'file_key': file_key
            }

        response = read_csv_handler(input)

        return response, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
