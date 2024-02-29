from rpc import RPC
import base64
import requests
import time
import random
import json
import os

def function_handler(context):
    if context["request_type"] == "GRPC":
        time.sleep(12)
        response = RPC(context, os.environ["FEATURE_STATUS"], [context["request"]])[0]
        return response, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
