from parliament import Context
from flask import Request
import base64
import requests
import json
import pprint
import random
import os
import datetime
import redis

if "LOGGING_NAME" in os.environ:
    loggingClient = redis.Redis(host=os.environ['LOGGING_URL'], password=os.environ['LOGGING_PASSWORD'])

pp = pprint.PrettyPrinter(indent=4)

def main(context: Context):
    if 'request' in context.keys():
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n")
        event = context.request.json

        try:
            pp
        except NameError:
            pp = pprint.PrettyPrinter(indent=4)

        bucket_name = event['Records'][0]['s3']['bucket']['name']
        file_key = event['Records'][0]['s3']['object']['key']

        input= {
                'bucket_name': bucket_name,
                'file_key': file_key
            }
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n")
        response = requests.get(url=os.environ["SENTIMENT_READ_CSV"], json=input)
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "2" + "\n")
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
