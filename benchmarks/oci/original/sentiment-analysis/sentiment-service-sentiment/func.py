import base64
import requests
import os
import json
import random
import datetime

def function_handler(context):
    if context["is_json"]:
        workflow_id = context["request"]["workflow_id"]
        workflow_depth = context["request"]["workflow_depth"]
        workflow_width = context["request"]["workflow_width"]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "3" + "\n", flush=True)
        event = context["request"]
        event["workflow_depth"] += 1

        feedback = event['feedback']
        response = {"polarity": -0.66}
        if response['polarity'] > 0.5:
            sentiment = "POSITIVE"
        elif response['polarity'] < -0.5:
            sentiment = "NEGATIVE"
        else:
            sentiment = "NEUTRAL"
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "4" + "\n", flush=True)
        response = requests.post(url=os.environ["SENTIMENT_SERVICE_RESULT"], json={
            'sentiment': sentiment,
            'reviewType': event['reviewType'],
            'reviewID': event['reviewID'],
            'customerID': event['customerID'],
            'productID': event['productID'],
            'feedback': event['feedback'], "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0
        })
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "5" + "\n", flush=True)
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
