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

if "LOGGING_NAME" in os.environ:
    loggingClient = redis.Redis(host=os.environ['LOGGING_URL'], password=os.environ['LOGGING_PASSWORD'])

pp = pprint.PrettyPrinter(indent=4)


def sfail_handler(req):
    return "SentimentFail: Fail: \"Sentiment Analysis Failed!\""

def main(context: Context):
    if 'request' in context.keys():
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n")
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
            if "LOGGING_NAME" in os.environ:
                loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n")
            response = requests.get(url=os.environ["SENTIMENT_DB_S"], json={
                'sentiment': sentiment,
                'reviewType': event['reviewType'],
                'reviewID': event['reviewID'],
                'customerID': event['customerID'],
                'productID': event['productID'],
                'feedback': event['feedback']
            })
            if "LOGGING_NAME" in os.environ:
                loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "2" + "\n")
            return response.text, 200
        else:
            return sfail_handler(event), 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
