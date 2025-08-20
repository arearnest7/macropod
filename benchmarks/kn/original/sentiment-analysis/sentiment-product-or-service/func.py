from parliament import Context
from flask import Request
import base64
import requests
import json
import random
import os
import datetime
import redis


def main(context: Context):
    if 'request' in context.keys():
        event = context.request.json

        results = ""
        if event["reviewType"] == "Product":
            results = requests.post(url=os.environ["SENTIMENT_PRODUCT_SENTIMENT"], json=event)
        elif event["reviewType"] == "Service":
            results = requests.post(url=os.environ["SENTIMENT_SERVICE_SENTIMENT"], json=event)
        else:
            results = requests.post(url=os.environ["SENTIMENT_CFAIL"], json=event)
        return results.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
