from parliament import Context
from flask import Request
import base64
import requests
import os
import redis
import random
import json

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def main(context: Context):
    if 'request' in context.keys():
        params = context.request.json
        num_of_file = int(params['num_of_file'])
        bucket = params['input_bucket']
        all_keys = []

        for key in ["reviews100mb.txt", "reviews10mb.txt", "reviews20mb.txt", "reviews50mb.txt"]:
            all_keys.append(key)
        print("Number of File : " + str(len(all_keys)))

        if num_of_file == len(all_keys):
            return requests.get(url=os.environ["FEATURE_REDUCER"], json=params).text, 200
        else:
            return requests.get(url=os.environ["FEATURE_WAIT"], json=params).text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
