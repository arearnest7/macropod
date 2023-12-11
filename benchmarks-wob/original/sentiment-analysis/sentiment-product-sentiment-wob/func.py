from parliament import Context
from flask import Request
import base64
import requests
import os
import json
import random

def main(context: Context):
    if 'request' in context.keys():
        event = context.request.json

        feedback = event['feedback']
        response = {"polarity": -0.66}
        if response['polarity'] > 0.5:
            sentiment = "POSITIVE"
        elif response['polarity'] < -0.5:
            sentiment = "NEGATIVE"
        else:
            sentiment = "NEUTRAL"
        response = requests.get(url=os.environ["SENTIMENT_PRODUCT_RESULT"], json={
            'sentiment': sentiment,
            'reviewType': event['reviewType'],
            'reviewID': event['reviewID'],
            'customerID': event['customerID'],
            'productID': event['productID'],
            'feedback': event['feedback']
        })
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
