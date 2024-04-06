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
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "0" + "\n", flush=True)
        event = context.request.json

        results = ""
        if event["reviewType"] == "Product":
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "1" + "\n", flush=True)
            results = requests.get(url=os.environ["SENTIMENT_PRODUCT_SENTIMENT"], json=event)
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "2" + "\n", flush=True)
        elif event["reviewType"] == "Service":
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "3" + "\n", flush=True)
            results = requests.get(url=os.environ["SENTIMENT_SERVICE_SENTIMENT"], json=event)
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "4" + "\n", flush=True)
        else:
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "5" + "\n", flush=True)
            results = requests.get(url=os.environ["SENTIMENT_CFAIL"], json=event)
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "6" + "\n", flush=True)
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "7" + "\n", flush=True)
        return results.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
