import base64
import requests
import os
import redis
import random
import json
import datetime

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def function_handler(context):
    if context["is_json"]:
        workflow_id = context["request"]["workflow_id"]
        workflow_depth = context["request"]["workflow_depth"]
        workflow_width = context["request"]["workflow_width"]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "3" + "\n", flush=True)
        params = context["request"]
        params["workflow_depth"] += 1
        num_of_file = int(params['num_of_file'])
        bucket = params['input_bucket']
        all_keys = []

        for key in ["reviews100mb.txt", "reviews10mb.txt", "reviews20mb.txt", "reviews50mb.txt"]:
            all_keys.append(key)
        print("Number of File : " + str(len(all_keys)))

        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "4" + "\n", flush=True)
        if num_of_file == len(all_keys):
            return requests.post(url=os.environ["FEATURE_REDUCER"], json=params).text, 200
        else:
            return requests.post(url=os.environ["FEATURE_WAIT"], json=params).text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
