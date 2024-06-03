from parliament import Context
from flask import Request
import base64
import requests
import json
import random
import os
import datetime
import redis


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
        #response = requests.post(url = 'http://' + OF_Gateway_IP + ':' + OF_Gateway_Port + '/function/sha>
        #    "Subject": 'Negative Review Received',
        #    "Message": 'Review (ID = %i) of %s (ID = %i) received with negative results from sentiment a>
        #    event['reviewType'], int(event['productID']), int(event['customerID']), event['feedback'])
        #})

        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "1" + "\n", flush=True)
        response = requests.post(url=os.environ["SENTIMENT_DB"], json={
            'sentiment': event['sentiment'],
            'reviewType': event['reviewType'],
            'reviewID': event['reviewID'],
            'customerID': event['customerID'],
            'productID': event['productID'],
            'feedback': event['feedback'], "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0
        })
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "2" + "\n", flush=True)
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
