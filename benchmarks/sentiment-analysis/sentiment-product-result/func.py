from rpc import Invoke_JSON
import base64
import json
import random
import os

def FunctionHandler(context):
    event = context["JSON"]
    results = ""
    if event["sentiment"] == "POSITIVE" or event["sentiment"] == "NEUTRAL":
        results = Invoke_JSON(context, "SENTIMENT_DB", event)
    elif event["sentiment"] == "NEGATIVE":
        results = Invoke_JSON(context, "SENTIMENT_SNS", event)
    else:
        results = Invoke_JSON(context, "SENTIMENT_SFAIL", event)

    return results[0], 200
