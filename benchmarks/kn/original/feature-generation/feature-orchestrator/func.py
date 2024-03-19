from parliament import Context
from flask import Request
import base64
import requests
import datetime
import redis
import json
from functools import partial
from multiprocessing.dummy import Pool as ThreadPool
import os
import random

redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

if "LOGGING_NAME" in os.environ:
    loggingClient = redis.Redis(host=os.environ['LOGGING_URL'], password=os.environ['LOGGING_PASSWORD'])

def invoke_lambda(bucket, dest, key):
    if "LOGGING_NAME" in os.environ:
        loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "3" + "\n")
    requests.post(url=os.environ["FEATURE_EXTRACTOR"], json={"input_bucket": bucket, "key": key.decode(), "dest": dest})
    if "LOGGING_NAME" in os.environ:
        loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "4" + "\n")

def main(context: Context):
    if 'request' in context.keys():
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n")
        params = context.request.json
        bucket = params['bucket']
        dest = str(random.randint(0, 10000000)) + "-" + bucket
        all_keys = []

        for key in redisClient.scan_iter(bucket + "-*"):
            all_keys.append(key)
        print("Number of File : " + str(len(all_keys)))
        print("File : " + str(all_keys))

        pool = ThreadPool(len(all_keys))
        pool.map(partial(invoke_lambda, bucket, dest), all_keys)
        pool.close()
        pool.join()

        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n")
        ret = requests.post(url=os.environ["FEATURE_WAIT"], json={"num_of_file": str(len(all_keys)), "input_bucket": dest}).text
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "2" + "\n")
        return ret, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
