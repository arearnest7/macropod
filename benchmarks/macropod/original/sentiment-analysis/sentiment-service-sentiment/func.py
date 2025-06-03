from rpc import Invoke_JSON
import base64
import requests
import os
import json
import random

def FunctionHandler(context):
    event = context["JSON"]

    feedback = event['feedback']
    response = {"polarity": -0.66}
    if "Bad" in feedback:
        response = {"polarity": -0.66}
    elif "Good" in feedback:
        response = {"polarity": 0.66}
    else:
        response = {"polarity": 0}
    if response['polarity'] > 0.5:
        sentiment = "POSITIVE"
    elif response['polarity'] < -0.5:
        sentiment = "NEGATIVE"
    else:
        sentiment = "NEUTRAL"
    payload = {
        'sentiment': sentiment,
        'reviewType': event['reviewType'],
        'reviewID': event['reviewID'],
        'customerID': event['customerID'],
        'productID': event['productID'],
        'feedback': event['feedback']
    }
    response = Invoke_JSON(context, "SENTIMENT_SERVICE_RESULT", payload)[0]
    return response.text, 200
