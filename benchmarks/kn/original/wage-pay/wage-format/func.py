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
        params = context.request.json
        params['INSURANCE'] = INSURANCE

        total = INSURANCE + params['base'] + params['merit']
        params['total'] = total

        realpay = (1-TAX) * (params['base'] + params['merit'])
        params['realpay'] = realpay

        response = requests.post(url=os.environ["WAGE_WRITE_RAW"], json=params)
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
