from rpc import RPC
import base64
import requests
import json
import random
import os

def function_handler(context):
    if context["request_type"] == "GRPC":
        event = json.loads(context["request"])
        results = ""
        if event["sentiment"] == "POSITIVE" or event["sentiment"] == "NEUTRAL":
            results = RPC(context, os.environ["SENTIMENT_DB"], [context["request"]])[0]
        elif event["sentiment"] == "NEGATIVE":
            results = RPC(context, os.environ["SENTIMENT_SNS"], [context["request"]])[0]
        else:
            results = RPC(context, os.environ["SENTIMENT_SFAIL"], [context["request"]])[0]

        return results, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
