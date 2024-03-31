from parliament import Context
from flask import Request
import base64

import requests
import re
import time
from sklearn.feature_extraction.text import TfidfVectorizer
import datetime
import redis
from functools import partial
from multiprocessing.dummy import Pool as ThreadPool
import os
import random
import json

redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

if "LOGGING_NAME" in os.environ:
    loggingClient = redis.Redis(host=os.environ['LOGGING_IP'], password=os.environ['LOGGING_PASSWORD'])

cleanup_re = re.compile('[^a-z]+')

def reducer_handler(req):
    params = req
    bucket = params['input_bucket']

    result = []
    latency = 0

    for key in redisClient.scan_iter(bucket + "-*"):
        body = redisClient.get(key).decode()
        start = time.time()
        word = body.replace("'", '').split(',')
        result.extend(word)
        latency += time.time() - start

    print(len(result))

    tfidf_vect = TfidfVectorizer().fit(result)
    feature = str(tfidf_vect.get_feature_names_out())
    feature = feature.lstrip('[').rstrip(']').replace(' ' , '')

    feature_key = 'feature.txt'
    redisClient.set(bucket + "-" + feature_key, str(feature))

    return str(latency)

def status_handler(req):
    params = req
    num_of_file = int(params['num_of_file'])
    bucket = params['input_bucket']
    all_keys = []

    for key in redisClient.scan_iter(bucket + "-*"):
        all_keys.append(key)
    print("Number of File : " + str(len(all_keys)))

    if num_of_file == len(all_keys):
        return reducer_handler(req)
    else:
        return wait_handler(req)

def wait_handler(req):
    time.sleep(12)
    response = status_handler(req)
    return response

def cleanup(sentence):
    sentence = sentence.lower()
    sentence = cleanup_re.sub(' ', sentence).strip()
    return sentence

def invoke_lambda(bucket, dest, key):
    if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "3" + "\n")
    requests.post(url=os.environ["FEATURE_EXTRACTOR_PARTIAL"], json={"input_bucket": bucket, "key": key.decode(), "dest": dest})
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
        ret = wait_handler({"num_of_file": str(len(all_keys)), "input_bucket": dest})
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "2" + "\n")
        return ret, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
