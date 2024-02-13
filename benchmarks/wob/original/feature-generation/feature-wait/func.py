from parliament import Context
from flask import Request
import base64
import requests
import time
import random
import json
import os

def main(context: Context):
    if 'request' in context.keys():
        time.sleep(12)
        response = requests.get(url=os.environ["FEATURE_STATUS"], json=context.request.json)
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
