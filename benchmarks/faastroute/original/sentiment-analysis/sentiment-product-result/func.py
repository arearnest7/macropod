from rpc import RPC
import base64
import requests
import json
import random
import os

def function_handler(context):
    if context["is_json"]:
        event = context["request"]
        results = ""
        if event["sentiment"] == "POSITIVE" or event["sentiment"] == "NEUTRAL":
            results = requests.get(url=os.environ["SENTIMENT_DB"], json=event)
        elif event["sentiment"] == "NEGATIVE":
            results = requests.get(url=os.environ["SENTIMENT_SNS"], json=event)
        else:
            results = requests.get(url=os.environ["SENTIMENT_SFAIL"], json=event)

        return results.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
