from rpc import RPC
import base64
import redis
import time
from sklearn.feature_extraction.text import TfidfVectorizer
import os
import json

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def function_handler(context):
    if context["request_type"] == "GRPC":
        params = json.loads(context["request"])
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

        return str(latency), 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
