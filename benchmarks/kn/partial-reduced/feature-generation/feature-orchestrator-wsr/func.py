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

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])


cleanup_re = re.compile('[^a-z]+')

def reducer_handler(req):
    params = req
    bucket = params['input_bucket']

    result = []
    latency = 0

    for key in ["reviews100mb.txt", "reviews10mb.txt", "reviews20mb.txt", "reviews50mb.txt"]:
        body = open(key, 'r').read()
        start = time.time()
        word = body.replace("'", '').split(',')
        result.extend(word)
        latency += time.time() - start

    print(len(result))

    tfidf_vect = TfidfVectorizer().fit(result)
    feature = str(tfidf_vect.get_feature_names_out())
    feature = feature.lstrip('[').rstrip(']').replace(' ' , '')

    feature_key = 'feature.txt'
    #redisClient.set(bucket + "-" + feature_key, str(feature))

    return str(latency)

def status_handler(req):
    params = req
    num_of_file = int(params['num_of_file'])
    bucket = params['input_bucket']
    all_keys = []

    for key in ["reviews100mb.txt", "reviews10mb.txt", "reviews20mb.txt", "reviews50mb.txt"]:
        all_keys.append(key)
    print("Number of File : " + str(len(all_keys)))

    if num_of_file == len(all_keys):
        return reducer_handler(req)
    else:
        return wait_handler(req)

def wait_handler(req):
    time.sleep(1)
    response = status_handler(req)
    return response

def cleanup(sentence):
    sentence = sentence.lower()
    sentence = cleanup_re.sub(' ', sentence).strip()
    return sentence

def invoke_lambda(bucket, dest, workflow_id, workflow_depth, key):
    requests.post(url=os.environ["FEATURE_EXTRACTOR_PARTIAL"], json={"input_bucket": bucket, "key": key, "dest": dest, "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0})

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
        params = context.request.json
        bucket = params['bucket']
        dest = str(random.randint(0, 10000000)) + "-" + bucket
        all_keys = []

        for key in ["reviews100mb.csv", "reviews10mb.csv", "reviews20mb.csv", "reviews50mb.csv"]:
            all_keys.append(key)
        print("Number of File : " + str(len(all_keys)))
        print("File : " + str(all_keys))

        pool = ThreadPool(len(all_keys))
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "1" + "\n", flush=True)
        pool.map(partial(invoke_lambda, bucket, dest, workflow_id, workflow_depth), all_keys)
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "2" + "\n", flush=True)
        pool.close()
        pool.join()

        ret = wait_handler({"num_of_file": str(len(all_keys)), "input_bucket": dest})
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "3" + "\n", flush=True)
        return ret, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
