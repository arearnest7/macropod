import base64
import requests
import json
import random
import os
import datetime

TAX = 0.0387
ROLES = ['staff', 'teamleader', 'manager']

def function_handler(context):
    if context["is_json"]:
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "10" + "\n", flush=True)
        params = context["request"]
        meritp = {'staff': 0, 'teamleader': 0, 'manager': 0}
        for role in ROLES:
            num = params['total']['statistics'][role+'-number']
            if num != 0:
                base = params['base']['statistics'][role]
                merit = params['merit']['statistics'][role]
                meritp[role] = merit / base
        params['statistics']['average-merit-percent'] = meritp
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "11" + "\n", flush=True)
        response = requests.get(url=os.environ["WAGE_WRITE_MERIT"], json={'id': params['id'], 'statistics': params['statistics'], 'operator' : params['operator']})
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "12" + "\n", flush=True)
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
