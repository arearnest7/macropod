from rpc import RPC
import base64
import requests
import time
import random
import json
import os

def FunctionHandler(context):
    if context["InvokeType"] == "GRPC":
        time.sleep(1)
        response = RPC(context, os.environ["FEATURE_STATUS"], [context["Request"]])[0]
        return response, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
