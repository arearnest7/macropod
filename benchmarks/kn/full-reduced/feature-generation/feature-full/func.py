from parliament import Context
from flask import Request
import base64
import requests
import redis
import json
from functools import partial
from multiprocessing.dummy import Pool as ThreadPool
import os
import random
import pandas as pd
import time
import re
from sklearn.feature_extraction.text import TfidfVectorizer

redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

cleanup_re = re.compile('[^a-z]+')

def reducer_handler(req):
    bucket = req['input_bucket']

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
    num_of_file = int(req['num_of_file'])
    bucket = req['input_bucket']
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

def extractor_handler(req):
    bucket = req['input_bucket']
    key = req['key']
    dest = req['dest']
    with open("/tmp/" + key, "w") as f:
        f.write(redisClient.get(key).decode())
    f.close()
    df = pd.read_csv("/tmp/" + key)

    start = time.time()
    df['Text'] = df['Text'].apply(cleanup)
    text = df['Text'].tolist()
    result = set()
    for item in text:
        result.update(item.split())
    print("Number of Feature : " + str(len(result)))

    feature = str(list(result))
    feature = feature.lstrip('[').rstrip(']').replace(' ', '')
    latency = time.time() - start
    print(latency)

    write_key = req['key'].split('.')[0] + ".txt"
    redisClient.set(dest + "-" + write_key, feature)
    return str(latency)

def invoke_lambda(bucket, dest, key):
    extractor_handler({"input_bucket": bucket, "key": key.decode(), "dest": dest})

def main(context: Context):
    if 'request' in context.keys():
        bucket = context.request.json['bucket']
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
        return wait_handler({"num_of_file": str(len(all_keys)), "input_bucket": dest}), 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
