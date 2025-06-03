from rpc import Invoke_JSON
import base64
import requests
import json
import random
import os

TAX = 0.0387
ROLES = ['staff', 'teamleader', 'manager']

def FunctionHandler(context):
    params = context["JSON"]

    realpay = {'staff': 0, 'teamleader': 0, 'manager': 0}
    for role in ROLES:
        num = params['total']['statistics'][role+'-number']
        if num != 0:
            base = params['base']['statistics'][role]
            merit = params['merit']['statistics'][role]
            realpay[role] = (1-TAX) * (base + merit) / num
    params['statistics']['average-realpay'] = realpay

    response = Invoke(context, "WAGE_MERIT", params])[0]
    return response, 200
