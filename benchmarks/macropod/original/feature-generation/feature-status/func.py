from rpc import RPC
import base64
import requests
import os
import redis
import random
import json

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def FunctionHandler(context):
    params = json.loads(context["Request"])
    num_of_file = int(params['num_of_file'])
    bucket = params['input_bucket']
    all_keys = []

    for key in ["reviews100mb.txt", "reviews10mb.txt", "reviews20mb.txt", "reviews50mb.txt"]:
        all_keys.append(key)
    print("Number of File : " + str(len(all_keys)))

    payload = []
    payload.append(context["Request"])
    if num_of_file == len(all_keys):
        return RPC(context, os.environ["FEATURE_REDUCER"], payload)[0], 200
    else:
        return RPC(context, os.environ["FEATURE_WAIT"], payload)[0], 200

