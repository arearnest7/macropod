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
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "10" + "\n", flush=True)
        params = context["request"]
        num_of_file = int(params['num_of_file'])
        bucket = params['input_bucket']
        all_keys = []

        for key in ["reviews100mb.txt", "reviews10mb.txt", "reviews20mb.txt", "reviews50mb.txt"]:
            all_keys.append(key)
        print("Number of File : " + str(len(all_keys)))

        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "11" + "\n", flush=True)
        if num_of_file == len(all_keys):
            return requests.get(url=os.environ["FEATURE_REDUCER"], json=params).text, 200
        else:
            return requests.get(url=os.environ["FEATURE_WAIT"], json=params).text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
