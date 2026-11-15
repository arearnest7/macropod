from rpc import Invoke_JSON
import base64
import json
import random
import os

def FunctionHandler(context):
    event = context["JSON"]

    results = ""
    if event["reviewType"] == "Product":
        results = Invoke_JSON(context, "SENTIMENT_PRODUCT_SENTIMENT", event)
    elif event["reviewType"] == "Service":
        results = Invoke_JSON(context, "SENTIMENT_SERVICE_SENTIMENT", event)
    else:
        results = Invoke_JSON(context, "SENTIMENT_CFAIL", event)
    return results[0], 200
