from parliament import Context
from flask import Request
import base64
import requests
import json
import time
from concurrent.futures import ThreadPoolExecutor
import os
import datetime
import redis
import random

TAX = 0.0387
INSURANCE = 1500
ROLES = ['staff', 'teamleader', 'manager']

#redisClient = redis.Redis(host=os.environ["REDIS_URL"],password=os.environ["REDIS_PASSWORD"])


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

def wage_sum_handler(req):
    params = json.loads(req)
    temp = json.loads(open(params["operator"], 'r').read())
    params["operator"] = temp["operator"]
    params["id"] = temp["id"]
    stats = {'total': params['total']['statistics']['total'] }
    params['statistics'] = stats

    return wage_avg_handler(params)

def stats_handler(req):
    manifest = []

    total = {'statistics': {'total': 0, 'staff-number': 0, 'teamleader-number': 0, 'manager-number': 0}}
    base = {'statistics': {'staff': 0, 'teamleader': 0, 'manager': 0}}
    merit = {'statistics': {'staff': 0, 'teamleader': 0, 'manager': 0}}

    for key in range(0,100):
        doc = json.loads(open(str(key), 'r').read())
        total['statistics']['total'] += doc['total']
        total['statistics'][doc['role']+'-number'] += 1
        base['statistics'][doc['role']] += doc['base']
        merit['statistics'][doc['role']] += doc['merit']
        manifest.append(str(key))

    fs = []
    with ThreadPoolExecutor(max_workers=len(manifest)) as executor:
        for obj in manifest:
            if obj != "raw/":
                fs.append(executor.submit(wage_sum_handler, {'total': total, 'base': base, 'merit': merit, 'operator': obj}))
    results = [f for f in fs]
    return "processed batch at " + str(time.time())

def write_raw_handler(req):
    params = req
    #redisClient.set("raw-" + str(params["id"]), json.dumps(req))
    response = requests.get(url=os.environ["WAGE_FULL"], json={})
    return response.text

def format_handler(req):
    params = req
    params['INSURANCE'] = INSURANCE

    total = INSURANCE + params['base'] + params['merit']
    params['total'] = total

    realpay = (1-TAX) * (params['base'] + params['merit'])
    params['realpay'] = realpay

    return write_raw_handler(params)

def validator_handler(req):
    event = req
    for param in ['id', 'name', 'role', 'base', 'merit', 'operator']:
        if param in ['name', 'role']:
            if not isinstance(event[param], str):
                return "fail: illegal params: " + str(event[param]) + " not string"
            elif param == 'role' and event[param] not in ROLES:
                return "fail: invalid role: " + str(event[param])
        elif param in ['id', 'operator']:
            if not isinstance(event[param], int):
                return "fail: illegal params: " + str(event[param]) + " not integer"
        elif param in ['base', 'merit']:
            if not isinstance(event[param], float):
                return "fail: illegal params: " + str(event[param]) + " not float"
            elif event[param] < 1 or event[param] > 8:
                return "fail: illegal params: " + str(event[param]) + " not between 1 and 8 inclusively"
        else:
            return "fail: missing param: " + param
    return format_handler(event)

def main(context: Context):
    if 'request' in context.keys():
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "0" + "\n", flush=True)
        params = context.request.json
        response = ""
        if len(params) > 0:
            response = validator_handler(params)
        else:
            response = stats_handler(params)
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "1" + "\n", flush=True)
        return response, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
