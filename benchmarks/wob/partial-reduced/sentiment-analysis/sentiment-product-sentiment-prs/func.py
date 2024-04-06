from parliament import Context
from flask import Request
import base64
import requests
import json
import pprint
import csv
import os
import random
import datetime
import redis


pp = pprint.PrettyPrinter(indent=4)


def sfail_handler(req):
    return "SentimentFail: Fail: \"Sentiment Analysis Failed!\""

def main(context: Context):
    if 'request' in context.keys():
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "0" + "\n", flush=True)
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
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "1" + "\n", flush=True)
            response = requests.get(url=os.environ["SENTIMENT_DB_S"], json={
                'sentiment': sentiment,
                'reviewType': event['reviewType'],
                'reviewID': event['reviewID'],
                'customerID': event['customerID'],
                'productID': event['productID'],
                'feedback': event['feedback']
            })
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "2" + "\n", flush=True)
            return response.text, 200
        else:
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "3" + "\n", flush=True)
            return sfail_handler(event), 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
