from parliament import Context
from flask import Request
import base64
import requests
import json
import random
import os
import datetime
import redis

if "LOGGING_NAME" in os.environ:
    loggingClient = redis.Redis(host=os.environ['LOGGING_IP'], password=os.environ['LOGGING_PASSWORD'])

def main(context: Context):
    if 'request' in context.keys():
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n")
        event = context.request.json

        results = ""
        if event["reviewType"] == "Product":
            if "LOGGING_NAME" in os.environ:
                loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n")
            results = requests.get(url=os.environ["SENTIMENT_PRODUCT_SENTIMENT"], json=event)
            if "LOGGING_NAME" in os.environ:
                loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "2" + "\n")
        elif event["reviewType"] == "Service":
            if "LOGGING_NAME" in os.environ:
                loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "3" + "\n")
            results = requests.get(url=os.environ["SENTIMENT_SERVICE_SENTIMENT"], json=event)
            if "LOGGING_NAME" in os.environ:
                loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "4" + "\n")
        else:
            if "LOGGING_NAME" in os.environ:
                loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "5" + "\n")
            results = requests.get(url=os.environ["SENTIMENT_CFAIL"], json=event)
            if "LOGGING_NAME" in os.environ:
                loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "6" + "\n")
        return results.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
