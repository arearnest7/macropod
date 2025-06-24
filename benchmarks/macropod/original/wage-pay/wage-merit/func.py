from rpc import Invoke_JSON
import base64
import json
import random
import os

TAX = 0.0387
ROLES = ['staff', 'teamleader', 'manager']

def FunctionHandler(context):
    params = context["JSON"]
    meritp = {'staff': 0, 'teamleader': 0, 'manager': 0}
    for role in ROLES:
        num = params['total']['statistics'][role+'-number']
        if num != 0:
            base = params['base']['statistics'][role]
            merit = params['merit']['statistics'][role]
            meritp[role] = merit / base
    params['statistics']['average-merit-percent'] = meritp
    response = Invoke_JSON(context, "WAGE_WRITE_MERIT", {'id': params['id'], 'statistics': params['statistics'], 'operator' : params['operator']})[0]
    return response, 200
