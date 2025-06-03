from rpc import Invoke_JSON
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
    temp = json.loads(open(params["operator"], 'r').read())
    params["operator"] = temp["operator"]
    params["id"] = temp["id"]
    stats = {'total': params['total']['statistics']['total'] }
    params['statistics'] = stats

    response = Invoke_JSON(context, "WAGE_AVG", params)[0]
    return response, 200
