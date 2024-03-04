from rpc import RPC
import base64
import requests
import json
import random
import os

TAX = 0.0387
ROLES = ['staff', 'teamleader', 'manager']

def function_handler(context):
    if context["InvokeType"] == "GRPC":
        params = json.loads(context["Request"])

        realpay = {'staff': 0, 'teamleader': 0, 'manager': 0}
        for role in ROLES:
            num = params['total']['statistics'][role+'-number']
            if num != 0:
                base = params['base']['statistics'][role]
                merit = params['merit']['statistics'][role]
                realpay[role] = (1-TAX) * (base + merit) / num
        params['statistics']['average-realpay'] = realpay

        response = RPC(context, os.environ["WAGE_MERIT"], [json.dumps(params).encode()])[0]
        return response, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
