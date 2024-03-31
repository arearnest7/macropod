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

if "LOGGING_NAME" in os.environ:
    loggingClient = redis.Redis(host=os.environ['LOGGING_IP'], password=os.environ['LOGGING_PASSWORD'])

def write_raw_handler(req):
    params = req
    redisClient.set("raw-" + str(params["id"]), json.dumps(req))
    if "LOGGING_NAME" in os.environ:
        loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "2" + "\n")
    response = requests.get(url=os.environ["WAGE_STATS_PARTIAL"], json={})
    if "LOGGING_NAME" in os.environ:
        loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "3" + "\n")
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
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n")
        event = context.request.json
        for param in ['id', 'name', 'role', 'base', 'merit', 'operator']:
            if param in ['name', 'role']:
                if not isinstance(event[param], str):
                    return "fail: illegal params: " + str(event[param]) + " not string", 200
                elif param == 'role' and event[param] not in ROLES:
                    return "fail: invalid role: " + str(event[param]), 200
            elif param in ['id', 'operator']:
                if not isinstance(event[param], int):
                    return "fail: illegal params: " + str(event[param]) + " not integer", 200
            elif param in ['base', 'merit']:
                if not isinstance(event[param], float):
                    return "fail: illegal params: " + str(event[param]) + " not float", 200
                elif event[param] < 1 or event[param] > 8:
                    return "fail: illegal params: " + str(event[param]) + " not between 1 and 8 inclusively", 200
            else:
                return "fail: missing param: " + param, 200
        ret = format_handler(event)
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n")
        return ret, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
