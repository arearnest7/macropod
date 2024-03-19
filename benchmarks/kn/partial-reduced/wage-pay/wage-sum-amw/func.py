from parliament import Context
from flask import Request
import base64
import requests
import json
import os
import sys
import datetime
import redis
import random

TAX = 0.0387
ROLES = ['staff', 'teamleader', 'manager']

redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

if "LOGGING_NAME" in os.environ:
    loggingClient = redis.Redis(host=os.environ['LOGGING_URL'], password=os.environ['LOGGING_PASSWORD'])

def write_merit_handler(req):
    params = req
    redisClient.set("merit-" + str(params["id"]), json.dumps(req))

    return str(params["id"]) + " statistics uploaded/updated"

def wage_merit_handler(req):
    params = req
    meritp = {'staff': 0, 'teamleader': 0, 'manager': 0}
    for role in ROLES:
        num = params['total']['statistics'][role+'-number']
        if num != 0:
            base = params['base']['statistics'][role]
            merit = params['merit']['statistics'][role]
            meritp[role] = merit / base
    params['statistics']['average-merit-percent'] = meritp
    return write_merit_handler({'id': params['id'], 'statistics': params['statistics'], 'operator' : params['operator']})

def wage_avg_handler(req):
    params = req

    realpay = {'staff': 0, 'teamleader': 0, 'manager': 0}
    for role in ROLES:
        num = params['total']['statistics'][role+'-number']
        if num != 0:
            base = params['base']['statistics'][role]
            merit = params['merit']['statistics'][role]
            realpay[role] = (1-TAX) * (base + merit) / num
    params['statistics']['average-realpay'] = realpay

    return wage_merit_handler(params)

def main(context: Context):
    if 'request' in context.keys():
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n")
        params = context.request.json
        temp = json.loads(redisClient.get(params["operator"]))
        params["operator"] = temp["operator"]
        params["id"] = temp["id"]
        stats = {'total': params['total']['statistics']['total'] }
        params['statistics'] = stats

        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n")
        return wage_avg_handler(params), 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
