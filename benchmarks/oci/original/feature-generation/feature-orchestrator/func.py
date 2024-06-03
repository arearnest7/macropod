import base64
import requests
import redis
import json
from functools import partial
from multiprocessing.dummy import Pool as ThreadPool
import os
import random
import datetime

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def invoke_lambda(bucket, dest, workflow_id, workflow_depth, key):
    requests.post(url=os.environ["FEATURE_EXTRACTOR"], json={"input_bucket": bucket, "key": key, "dest": dest, "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0})

def function_handler(context):
    if context["is_json"]:
        workflow_id = context["request"]["workflow_id"]
        workflow_depth = context["request"]["workflow_depth"]
        workflow_width = context["request"]["workflow_width"]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "3" + "\n", flush=True)
        params = context["request"]
        bucket = params['bucket']
        dest = str(random.randint(0, 10000000)) + "-" + bucket
        all_keys = []

        for key in ["reviews100mb.csv", "reviews10mb.csv", "reviews20mb.csv", "reviews50mb.csv"]:
            all_keys.append(key)
        print("Number of File : " + str(len(all_keys)))
        print("File : " + str(all_keys))

        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "4" + "\n", flush=True)
        pool = ThreadPool(len(all_keys))
        pool.map(partial(invoke_lambda, bucket, dest, context["request"]["workflow_id"], context["request"]["workflow_depth"]), all_keys)
        pool.close()
        pool.join()
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "5" + "\n", flush=True)

        return requests.post(url=os.environ["FEATURE_WAIT"], json={"num_of_file": str(len(all_keys)), "input_bucket": dest, "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0}).text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
