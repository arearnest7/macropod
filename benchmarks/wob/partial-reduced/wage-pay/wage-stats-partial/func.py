from parliament import Context
from flask import Request
import base64
import requests
import json
import time
from concurrent.futures import ThreadPoolExecutor
import os
import sys
import datetime
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

if "LOGGING_NAME" in os.environ:
    loggingClient = redis.Redis(host=os.environ['LOGGING_IP'], password=os.environ['LOGGING_PASSWORD'])

def main(context: Context):
    if 'request' in context.keys():
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n")
        manifest = []

        total = {'statistics': {'total': 0, 'staff-number': 0, 'teamleader-number': 0, 'manager-number': 0}}
        base = {'statistics': {'staff': 0, 'teamleader': 0, 'manager': 0}}
        merit = {'statistics': {'staff': 0, 'teamleader': 0, 'manager': 0}}

        for key in range(0, 100):
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
                    fs.append(executor.submit(requests.get, url=os.environ["WAGE_SUM_AMW"], json={'total': total, 'base': base, 'merit': merit, 'operator': obj}))
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n")
        results = [f for f in fs]
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "2" + "\n")
        return "processed batch at " + str(time.time()), 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
