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


pp = pprint.PrettyPrinter(indent=4)
#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def cfail_handler(req):
    return "CategorizationFail: Fail: \"Input CSV could not be categorised into 'Product' or 'Service'.\""

def product_or_service_handler(req):
    workflow_id = req["workflow_id"]
    workflow_depth = req["workflow_depth"]
    workflow_width = req["workflow_width"]
    event = req

    results = ""
    if event["reviewType"] == "Product":
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "1" + "\n", flush=True)
        response = requests.post(url=os.environ["SENTIMENT_PRODUCT_SENTIMENT_PRS"], json=event)
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "2" + "\n", flush=True)
        results = response.text
    elif event["reviewType"] == "Service":
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "3" + "\n", flush=True)
        response = requests.post(url=os.environ["SENTIMENT_SERVICE_SENTIMENT_SRS"], json=event)
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "4" + "\n", flush=True)
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
        row["workflow_id"] = req["workflow_id"]
        row["workflow_depth"] = req["workflow_depth"] + 1
        row["workflow_width"] = 0
        return product_or_service_handler(row)

def main(context: Context):
    if 'request' in context.keys():
        workflow_id = str(random.randint(0, 10000000))
        workflow_depth = 0
        workflow_width = 0
        if "workflow_id" in context.request.json:
            workflow_id = context.request.json["workflow_id"]
            workflow_depth = context.request.json["workflow_depth"]
            workflow_width = context.request.json["workflow_width"]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "0" + "\n", flush=True)
        event = context.request.json

        try:
            pp
        except NameError:
            pp = pprint.PrettyPrinter(indent=4)

        bucket_name = event['Records'][0]['s3']['bucket']['name']
        file_key = event['Records'][0]['s3']['object']['key']


        input= {
                'bucket_name': bucket_name,
                'file_key': file_key,
                "workflow_id": workflow_id,
                "workflow_depth": workflow_depth + 1,
                "workflow_width": 0
            }

        response = read_csv_handler(input)
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "5" + "\n", flush=True)
        return response, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
