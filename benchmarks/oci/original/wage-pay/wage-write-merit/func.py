import base64
import json
import os
import sys
import redis
import random
import datetime

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def function_handler(context):
    if context["is_json"]:
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "10" + "\n", flush=True)
        params = context["request"]

        #redisClient.set("merit-" + str(params["id"]), json.dumps(params))

        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "11" + "\n", flush=True)
        return str(params["id"]) + " statistics uploaded/updated", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
