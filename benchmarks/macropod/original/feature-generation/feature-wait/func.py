from rpc import RPC
import base64
import requests
import time
import random
import json
import os

def FunctionHandler(context):
    time.sleep(1)
    payload = []
    payload.append(context["Request"])
    response = RPC(context, os.environ["FEATURE_STATUS"], payload)[0]
    return response, 200
