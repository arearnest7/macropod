from rpc import RPC
import base64
import redis
import pandas as pd
import time
import re
import os
import json

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

cleanup_re = re.compile('[^a-z]+')

def cleanup(sentence):
    sentence = sentence.lower()
    sentence = cleanup_re.sub(' ', sentence).strip()
    return sentence

def FunctionHandler(context):
    params = json.loads(context["Request"])
    bucket = params['input_bucket']
    key = params['key']
    dest = params['dest']
    with open("/tmp/" + key, "w") as f:
        f.write(open(key, 'r').read())
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

    write_key = params['key'].split('.')[0] + ".txt"
    #redisClient.set(dest + "-" + write_key, feature)
    return str(latency), 200
