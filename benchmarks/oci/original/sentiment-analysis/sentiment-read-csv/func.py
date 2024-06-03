import base64
import requests
import csv
import json
import os
import redis
import random
import datetime

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def function_handler(context):
    if context["is_json"]:
        workflow_id = context["request"]["workflow_id"]
        workflow_depth = context["request"]["workflow_depth"]
        workflow_width = context["request"]["workflow_width"]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "3" + "\n", flush=True)
        event = context["request"]
        event["workflow_depth"] += 1

        bucket_name = event['bucket_name']
        file_key = event['file_key']

        response = open(file_key, 'r').read()

        lines = response.split('\n')

        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "4" + "\n", flush=True)
        for row in csv.DictReader(lines):
            row["workflow_id"] = workflow_id
            row["workflow_depth"] = workflow_depth + 1
            row["workflow_width"] = 0
            response = requests.post(url=os.environ["SENTIMENT_PRODUCT_OR_SERVICE"], json=row)
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "5" + "\n", flush=True)
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
