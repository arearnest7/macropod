import base64
import requests
import json
import random
import os

TAX = 0.0387
INSURANCE = 1500

def function_handler(context):
    if context["is_json"]:
        params = context["request"]
        print(type(params))
        params['INSURANCE'] = INSURANCE

        total = INSURANCE + params['base'] + params['merit']
        params['total'] = total

        realpay = (1-TAX) * (params['base'] + params['merit'])
        params['realpay'] = realpay

        response = requests.get(url=os.environ["WAGE_WRITE_RAW"], json=params)
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
