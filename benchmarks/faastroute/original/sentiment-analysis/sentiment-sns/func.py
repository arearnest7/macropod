from rpc import RPC
import base64
import requests
import json
import random
import os

def function_handler(context):
    if context["is_json"]:
        event = context["request"]
        #response = requests.get(url = 'http://' + OF_Gateway_IP + ':' + OF_Gateway_Port + '/function/sha>
        #    "Subject": 'Negative Review Received',
        #    "Message": 'Review (ID = %i) of %s (ID = %i) received with negative results from sentiment a>
        #    event['reviewType'], int(event['productID']), int(event['customerID']), event['feedback'])
        #})

        response = requests.get(url=os.environ["SENTIMENT_DB"], json={
            'sentiment': event['sentiment'],
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
