from rpc import RPC
import base64
import requests
import json
import random
import os

TAX = 0.0387
ROLES = ['staff', 'teamleader', 'manager']

def function_handler(context):
    if context["request_type"] == "GRPC":
        params = json.loads(context["request"])
        meritp = {'staff': 0, 'teamleader': 0, 'manager': 0}
        for role in ROLES:
            num = params['total']['statistics'][role+'-number']
            if num != 0:
                base = params['base']['statistics'][role]
                merit = params['merit']['statistics'][role]
                meritp[role] = merit / base
        params['statistics']['average-merit-percent'] = meritp
        response = RPC(os.environ["WAGE_WRITE_MERIT"], [json.dumps({'id': params['id'], 'statistics': params['statistics'], 'operator' : params['operator']})], context["workflow_id"])[0]
        return response, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
