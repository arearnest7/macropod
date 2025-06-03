from rpc import Invoke_JSON
import base64
import requests
import json
import random
import os

def FunctionHandler(context):
    event = context["JSON"]

    results = ""
    if event["reviewType"] == "Product":
        results = Invoke_JSON(context, "SENTIMENT_PRODUCT_SENTIMENT", event)[0]
    elif event["reviewType"] == "Service":
        results = Invoke_JSON(context, "SENTIMENT_SERVICE_SENTIMENT", event)[0]
    else:
        results = Invoke_JSON(context, "SENTIMENT_CFAIL", event)[0]
    return results, 200
