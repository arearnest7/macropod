from rpc import RPC
import base64
import requests
import json
import random
import os

def FunctionHandler(context):
    event = json.loads(context["Request"])
    results = ""
    payload = []
    payload.append(context["Request"])
    if event["sentiment"] == "POSITIVE" or event["sentiment"] == "NEUTRAL":
        results = RPC(context, os.environ["SENTIMENT_DB"], payload)[0]
    elif event["sentiment"] == "NEGATIVE":
        results = RPC(context, os.environ["SENTIMENT_SNS"], payload)[0]
    else:
        results = RPC(context, os.environ["SENTIMENT_SFAIL"], payload)[0]

    return results, 200
