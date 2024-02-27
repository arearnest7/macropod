from rpc import RPC
import base64
import requests
import os
import json
import random

def function_handler(context):
    if context["request_type"] == "GRPC":
        event = json.loads(context["request"])

        feedback = event['feedback']
        response = {"polarity": -0.66}
        if response['polarity'] > 0.5:
            sentiment = "POSITIVE"
        elif response['polarity'] < -0.5:
            sentiment = "NEGATIVE"
        else:
            sentiment = "NEUTRAL"
        response = RPC(os.environ["SENTIMENT_SERVICE_RESULT"], [json.dumps({
            'sentiment': sentiment,
            'reviewType': event['reviewType'],
            'reviewID': event['reviewID'],
            'customerID': event['customerID'],
            'productID': event['productID'],
            'feedback': event['feedback']
        })], context["workflow_id"])
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200