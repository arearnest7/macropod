from rpc import RPC
import base64
import requests
import json
import random
import os

TAX = 0.0387
INSURANCE = 1500

def function_handler(context):
    if context["request_type"] == "GRPC":
        params = json.loads(context["request"])
        print(type(params))
        params['INSURANCE'] = INSURANCE

        total = INSURANCE + params['base'] + params['merit']
        params['total'] = total

        realpay = (1-TAX) * (params['base'] + params['merit'])
        params['realpay'] = realpay

        response = RPC(os.environ["WAGE_WRITE_RAW"], [json.dumps(params)], context["workflow_id"])[0]
        return response, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
