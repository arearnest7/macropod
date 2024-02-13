from parliament import Context
from flask import Request
import base64
import redis
import time
from sklearn.feature_extraction.text import TfidfVectorizer
import os
import json

redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def main(context: Context):
    if 'request' in context.keys():
        params = context.request.json
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

        return str(latency), 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
