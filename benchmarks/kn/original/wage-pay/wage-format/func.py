from parliament import Context
from flask import Request
import base64
import requests
import json
import random
import os
import datetime
import redis


TAX = 0.0387
INSURANCE = 1500

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
        params["workflow_depth"] += 1
        print(type(params))
        params['INSURANCE'] = INSURANCE

        total = INSURANCE + params['base'] + params['merit']
        params['total'] = total

        realpay = (1-TAX) * (params['base'] + params['merit'])
        params['realpay'] = realpay

        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "1" + "\n", flush=True)
        response = requests.post(url=os.environ["WAGE_WRITE_RAW"], json=params)
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "2" + "\n", flush=True)
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
