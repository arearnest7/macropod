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
ROLES = ['staff', 'teamleader', 'manager']

def main(context: Context):
    if 'request' in context.keys():
        params = context.request.json
        meritp = {'staff': 0, 'teamleader': 0, 'manager': 0}
        for role in ROLES:
            num = params['total']['statistics'][role+'-number']
            if num != 0:
                base = params['base']['statistics'][role]
                merit = params['merit']['statistics'][role]
                meritp[role] = merit / base
        params['statistics']['average-merit-percent'] = meritp
        response = requests.post(url=os.environ["WAGE_WRITE_MERIT"], json={'id': params['id'], 'statistics': params['statistics'], 'operator' : params['operator']})
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
