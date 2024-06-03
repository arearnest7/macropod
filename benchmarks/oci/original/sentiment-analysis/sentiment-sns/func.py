import base64
import requests
import json
import random
import os
import datetime

def function_handler(context):
    if context["is_json"]:
        workflow_id = context["request"]["workflow_id"]
        workflow_depth = context["request"]["workflow_depth"]
        workflow_width = context["request"]["workflow_width"]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "3" + "\n", flush=True)
        event = context["request"]
        event["workflow_depth"] += 1
        #response = requests.post(url = 'http://' + OF_Gateway_IP + ':' + OF_Gateway_Port + '/function/sha>
        #    "Subject": 'Negative Review Received',
        #    "Message": 'Review (ID = %i) of %s (ID = %i) received with negative results from sentiment a>
        #    event['reviewType'], int(event['productID']), int(event['customerID']), event['feedback'])
        #})

        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "4" + "\n", flush=True)
        response = requests.post(url=os.environ["SENTIMENT_DB"], json={
            'sentiment': event['sentiment'],
            'reviewType': event['reviewType'],
            'reviewID': event['reviewID'],
            'customerID': event['customerID'],
            'productID': event['productID'],
            'feedback': event['feedback']
        })
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "5" + "\n", flush=True)
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
