from rpc import Invoke
import base64
import requests
import json
import os
import sys
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def FunctionHandler(context):
    params = context["JSON"]
    #redisClient.set("raw-" + str(params["id"]), json.dumps(params))
    response = Invoke(context, "WAGE_STATS", "")[0]
    return response, 200
