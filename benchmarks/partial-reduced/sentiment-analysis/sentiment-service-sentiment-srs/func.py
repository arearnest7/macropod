from parliament import Context
from flask import Request
import base64
import requests
import json
import pprint
import csv
import os
import random

pp = pprint.PrettyPrinter(indent=4)


def sfail_handler(req):
    return "SentimentFail: Fail: \"Sentiment Analysis Failed!\""

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

        if sentiment in ["POSITIVE", "NEGATIVE", "NEUTRAL"]:
            response = requests.get(url=os.environ["SENTIMENT_DB_S"], json={
                'sentiment': sentiment,
                'reviewType': event['reviewType'],
                'reviewID': event['reviewID'],
                'customerID': event['customerID'],
                'productID': event['productID'],
                'feedback': event['feedback']
            })
            return response.text, 200
        else:
            return sfail_handler(event), 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
