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

pp = pprint.PrettyPrinter(indent=4)


def main(context: Context):
    if 'request' in context.keys():
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
        response = requests.post(url=os.environ["SENTIMENT_READ_CSV"], json=input)
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
