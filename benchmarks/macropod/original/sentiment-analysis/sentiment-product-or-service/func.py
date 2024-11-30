from rpc import RPC
import base64
import requests
import json
import random
import os

def FunctionHandler(context):
    event = json.loads(context["Request"])

    results = ""
    payload = []
    payload.append(context["Request"])
    if event["reviewType"] == "Product":
        results = RPC(context, os.environ["SENTIMENT_PRODUCT_SENTIMENT"], payload)[0]
    elif event["reviewType"] == "Service":
        results = RPC(context, os.environ["SENTIMENT_SERVICE_SENTIMENT"], payload)[0]
    else:
        results = RPC(context, os.environ["SENTIMENT_CFAIL"], payload)[0]
    return results, 200
