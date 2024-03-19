from parliament import Context
from flask import Request
import base64
import requests
import json
import random
import os
import datetime
import redis

if "LOGGING_NAME" in os.environ:
    loggingClient = redis.Redis(host=os.environ['LOGGING_URL'], password=os.environ['LOGGING_PASSWORD'])

TAX = 0.0387
ROLES = ['staff', 'teamleader', 'manager']

def main(context: Context):
    if 'request' in context.keys():
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n")
        params = context.request.json

        realpay = {'staff': 0, 'teamleader': 0, 'manager': 0}
        for role in ROLES:
            num = params['total']['statistics'][role+'-number']
            if num != 0:
                base = params['base']['statistics'][role]
                merit = params['merit']['statistics'][role]
                realpay[role] = (1-TAX) * (base + merit) / num
        params['statistics']['average-realpay'] = realpay

        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n")
        response = requests.get(url=os.environ["WAGE_MERIT"], json=params)
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "2" + "\n")
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
