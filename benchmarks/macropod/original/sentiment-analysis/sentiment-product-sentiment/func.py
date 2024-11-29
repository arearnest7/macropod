from rpc import RPC
import base64
import requests
import os
import json
import random

def FunctionHandler(context):
    event = json.loads(context["Request"])

    feedback = event['feedback']
    response = {"polarity": -0.66}
    if response['polarity'] > 0.5:
       sentiment = "POSITIVE"
    elif response['polarity'] < -0.5:
       sentiment = "NEGATIVE"
    else:
        sentiment = "NEUTRAL"
    response = RPC(context, os.environ["SENTIMENT_PRODUCT_RESULT"], [json.dumps({
        'sentiment': sentiment,
        'reviewType': event['reviewType'],
        'reviewID': event['reviewID'],
        'customerID': event['customerID'],
        'productID': event['productID'],
        'feedback': event['feedback']
    }).encode()])[0]
    return response, 200
