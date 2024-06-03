from rpc import RPC
import base64
import requests
import json
import time
from concurrent.futures import ThreadPoolExecutor
import os
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def FunctionHandler(context):
    if context["InvokeType"] == "GRPC":
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
                    # fs.append(executor.submit(requests.post, url=os.environ["WAGE_SUM"], json={'total': total, 'base': base, 'merit': merit, 'operator': obj}))
                    fs.append(json.dumps({'total': total, 'base': base, 'merit': merit, 'operator': obj}).encode())
        #results = [f for f in fs]
        results = RPC(context, os.environ["WAGE_SUM"], fs)
        return "processed batch at " + str(time.time()), 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
