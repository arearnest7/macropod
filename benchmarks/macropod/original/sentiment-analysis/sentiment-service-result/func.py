from rpc import Invoke_JSON
import base64
import requests
import json
import random
import os

def FunctionHandler(context):
    event = context["JSON"]
    results = ""
    if event["sentiment"] == "POSITIVE" or event["sentiment"] == "NEUTRAL":
        results = Invoke_JSON(context, "SENTIMENT_DB", event)[0]
    elif event["sentiment"] == "NEGATIVE":
        results = Invoke_JSON(context, "SENTIMENT_SNS", event)[0]
    else:
        results = Invoke_JSON(context, "SENTIMENT_SFAIL", event)[0]

    return results, 200
