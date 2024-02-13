import base64
import requests
import json
import random
import os

def function_handler(context):
    if context["is_json"]:
        event = context["request"]

        results = ""
        if event["reviewType"] == "Product":
            results = requests.get(url=os.environ["SENTIMENT_PRODUCT_SENTIMENT"], json=event)
        elif event["reviewType"] == "Service":
            results = requests.get(url=os.environ["SENTIMENT_SERVICE_SENTIMENT"], json=event)
        else:
            results = requests.get(url=os.environ["SENTIMENT_CFAIL"], json=event)
        return results.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
