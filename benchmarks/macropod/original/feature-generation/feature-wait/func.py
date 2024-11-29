from rpc import RPC
import base64
import requests
import time
import random
import json
import os

def FunctionHandler(context):
    time.sleep(1)
    response = RPC(context, os.environ["FEATURE_STATUS"], [context["Request"]])[0]
    return response, 200
