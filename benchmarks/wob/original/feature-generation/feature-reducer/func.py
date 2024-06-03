from parliament import Context
from flask import Request
import base64
import datetime
import redis
import time
from sklearn.feature_extraction.text import TfidfVectorizer
import os
import json
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])


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
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "1" + "\n", flush=True)
        return str(latency), 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
