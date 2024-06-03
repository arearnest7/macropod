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

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def write_merit_handler(req):
    params = req
    #redisClient.set("merit-" + str(params["id"]), json.dumps(req))

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
        workflow_id = str(random.randint(0, 10000000))
        workflow_depth = 0
        workflow_width = 0
        if "workflow_id" in context.request.json:
            workflow_id = context.request.json["workflow_id"]
            workflow_depth = context.request.json["workflow_depth"]
            workflow_width = context.request.json["workflow_width"]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "0" + "\n", flush=True)
        params = context.request.json
        temp = json.loads(open(params["operator"], 'r').read())
        params["operator"] = temp["operator"]
        params["id"] = temp["id"]
        stats = {'total': params['total']['statistics']['total'] }
        params['statistics'] = stats
        ret = wage_avg_handler(params)
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "1" + "\n", flush=True)
        return ret, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
