from rpc import RPC
import base64
import requests
import time
import random
import json
import os

def function_handler(context):
    if context["is_json"]:
        time.sleep(12)
        response = requests.get(url=os.environ["FEATURE_STATUS"], json=context["request"])
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
