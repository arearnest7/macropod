from parliament import Context
from flask import Request
import base64
import requests
import json
import pprint
import csv
import os
import datetime
import redis
import random

if "LOGGING_NAME" in os.environ:
    loggingClient = redis.Redis(host=os.environ['LOGGING_URL'], password=os.environ['LOGGING_PASSWORD'])

pp = pprint.PrettyPrinter(indent=4)
#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def cfail_handler(req):
    return "CategorizationFail: Fail: \"Input CSV could not be categorised into 'Product' or 'Service'.\""

def product_or_service_handler(req):
    event = req

    results = ""
    if event["reviewType"] == "Product":
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "2" + "\n")
        response = requests.get(url=os.environ["SENTIMENT_PRODUCT_SENTIMENT_PRS"], json=event)
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "3" + "\n")
        results = response.text
    elif event["reviewType"] == "Service":
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "4" + "\n")
        response = requests.get(url=os.environ["SENTIMENT_SERVICE_SENTIMENT_SRS"], json=event)
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "5" + "\n")
        results = response.text
    else:
        results = cfail_handler(event)
    return results

def read_csv_handler(req):
    event = req

    bucket_name = event['bucket_name']
    file_key = event['file_key']

    response = open(file_key, 'r').read()

    lines = response.split('\n')

    for row in csv.DictReader(lines):
        return product_or_service_handler(row)

def main(context: Context):
    if 'request' in context.keys():
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n")
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

        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n")
        return response, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
