from rpc import Invoke_JSON
import base64
import requests
import json
import random
import os

TAX = 0.0387
INSURANCE = 1500

def FunctionHandler(context):
    params = context["Text"]
    print(type(params))
    params['INSURANCE'] = INSURANCE

    total = INSURANCE + params['base'] + params['merit']
    params['total'] = total

    realpay = (1-TAX) * (params['base'] + params['merit'])
    params['realpay'] = realpay

    response = Invoke_JSON(context, "WAGE_WRITE_RAW", params)[0]
    return response, 200
