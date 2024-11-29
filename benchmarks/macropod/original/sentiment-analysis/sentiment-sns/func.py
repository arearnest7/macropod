from rpc import RPC
import base64
import requests
import json
import random
import os

def FunctionHandler(context):
    event = json.loads(context["Request"])
    #response = requests.post(url = 'http://' + OF_Gateway_IP + ':' + OF_Gateway_Port + '/function/sha>
    #    "Subject": 'Negative Review Received',
    #    "Message": 'Review (ID = %i) of %s (ID = %i) received with negative results from sentiment a>
    #    event['reviewType'], int(event['productID']), int(event['customerID']), event['feedback'])
    #})

    response = RPC(context, os.environ["SENTIMENT_DB"], [json.dumps({
        'sentiment': event['sentiment'],
        'reviewType': event['reviewType'],
        'reviewID': event['reviewID'],
        'customerID': event['customerID'],
        'productID': event['productID'],
        'feedback': event['feedback']
    }).encode()])[0]
    return response, 200
