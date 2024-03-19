from rpc import RPC
import base64
import requests
import json
import random
import os

def function_handler(context):
    if context["InvokeType"] == "GRPC":
        event = json.loads(context["Request"])
        results = ""
        if event["sentiment"] == "POSITIVE" or event["sentiment"] == "NEUTRAL":
            results = RPC(context, os.environ["SENTIMENT_DB"], [context["Request"]])[0]
        elif event["sentiment"] == "NEGATIVE":
            results = RPC(context, os.environ["SENTIMENT_SNS"], [context["Request"]])[0]
        else:
            results = RPC(context, os.environ["SENTIMENT_SFAIL"], [context["Request"]])[0]

        return results, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
