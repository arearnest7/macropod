from rpc import RPC
import base64
import requests
import json
import random
import os

def function_handler(context):
    if context["InvokeType"] == "GRPC":
        event = json.loads(context["Request"])

        results = ""
        if event["reviewType"] == "Product":
            results = RPC(context, os.environ["SENTIMENT_PRODUCT_SENTIMENT"], [context["Request"]])[0]
        elif event["reviewType"] == "Service":
            results = RPC(context, os.environ["SENTIMENT_SERVICE_SENTIMENT"], [context["Request"]])[0]
        else:
            results = RPC(context, os.environ["SENTIMENT_CFAIL"], [context["Request"]])[0]
        return results, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
