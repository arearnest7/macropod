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
        workflow_id = str(random.randint(0, 10000000))
        workflow_depth = 0
        workflow_width = 0
        if "workflow_id" in context.request.json:
            workflow_id = context.request.json["workflow_id"]
            workflow_depth = context.request.json["workflow_depth"]
            workflow_width = context.request.json["workflow_width"]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "0" + "\n", flush=True)
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
            print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "1" + "\n", flush=True)
            response = requests.post(url=os.environ["SENTIMENT_DB_S"], json={
                'sentiment': sentiment,
                'reviewType': event['reviewType'],
                'reviewID': event['reviewID'],
                'customerID': event['customerID'],
                'productID': event['productID'],
                'feedback': event['feedback'], "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0
            })
            print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "2" + "\n", flush=True)
            return response.text, 200
        else:
            print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "3" + "\n", flush=True)
            return sfail_handler(event), 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
