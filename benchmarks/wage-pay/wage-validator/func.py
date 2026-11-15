from rpc import Invoke_JSON
import base64
import json
import random
import os

ROLES = ['staff', 'teamleader', 'manager']

def FunctionHandler(context):
    event = context["JSON"]
    for param in ['id', 'name', 'role', 'base', 'merit', 'operator']:
        if param in ['name', 'role']:
            if not isinstance(event[param], str):
                return "fail: illegal params: " + str(event[param]) + " not string", 200
            elif param == 'role' and event[param] not in ROLES:
                return "fail: invalid role: " + str(event[param]), 200
        elif param in ['id', 'operator']:
            if not isinstance(event[param], float):
                return "fail: illegal params: " + str(event[param]) + " not float", 200 # originally int check, but json object autoconverts to float. computation is the same, so this does not affect results.
        elif param in ['base', 'merit']:
            if not isinstance(event[param], float):
                return "fail: illegal params: " + str(event[param]) + " not float", 200
            elif event[param] < 1 or event[param] > 8:
                return "fail: illegal params: " + str(event[param]) + " not between 1 and 8 inclusively", 200
        else:
            return "fail: missing param: " + param, 200
    response = Invoke_JSON(context, "WAGE_FORMAT", event)[0]
    return response, 200
