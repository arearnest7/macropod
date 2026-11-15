from rpc import Invoke
import base64
import json
import os
import sys
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def FunctionHandler(context):
    params = context["JSON"]

    #redisClient.set("merit-" + str(params["id"]), json.dumps(params))

    return str(params["id"]) + " statistics uploaded/updated", 200
