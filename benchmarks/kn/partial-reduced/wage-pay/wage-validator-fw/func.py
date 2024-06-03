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
INSURANCE = 1500
ROLES = ['staff', 'teamleader', 'manager']

redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])


def write_raw_handler(req):
    workflow_id = req["workflow_id"]
    workflow_depth = req["workflow_depth"]
    workflow_width = req["workflow_width"]
    params = req
    redisClient.set("raw-" + str(params["id"]), json.dumps(req))
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "7" + "\n", flush=True)
    response = requests.post(url=os.environ["WAGE_STATS_PARTIAL"], json={"workflow_id": req["workflow_id"], "workflow_depth": req["workflow_depth"] + 1, "workflow_width": 0})
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "8" + "\n", flush=True)
    return response.text

def format_handler(req):
    params = req
    params['INSURANCE'] = INSURANCE

    total = INSURANCE + params['base'] + params['merit']
    params['total'] = total

    realpay = (1-TAX) * (params['base'] + params['merit'])
    params['realpay'] = realpay

    return write_raw_handler(params)

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
        event = context.request.json
        event["workflow_id"] = workflow_id
        event["workflow_depth"] = workflow_depth + 1
        event["workflow_width"] = 0
        for param in ['id', 'name', 'role', 'base', 'merit', 'operator']:
            if param in ['name', 'role']:
                if not isinstance(event[param], str):
                    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "1" + "\n", flush=True)
                    return "fail: illegal params: " + str(event[param]) + " not string", 200
                elif param == 'role' and event[param] not in ROLES:
                    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "2" + "\n", flush=True)
                    return "fail: invalid role: " + str(event[param]), 200
            elif param in ['id', 'operator']:
                if not isinstance(event[param], int):
                    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "3" + "\n", flush=True)
                    return "fail: illegal params: " + str(event[param]) + " not integer", 200
            elif param in ['base', 'merit']:
                if not isinstance(event[param], float):
                    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "4" + "\n", flush=True)
                    return "fail: illegal params: " + str(event[param]) + " not float", 200
                elif event[param] < 1 or event[param] > 8:
                    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "5" + "\n", flush=True)
                    return "fail: illegal params: " + str(event[param]) + " not between 1 and 8 inclusively", 200
            else:
                print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "6" + "\n", flush=True)
                return "fail: missing param: " + param, 200
        ret = format_handler(event)
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "9" + "\n", flush=True)
        return ret, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
