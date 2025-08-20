from rpc import Invoke_JSON
import base64
import json
import random
import os

def FunctionHandler(context):
    event = context["JSON"]
    #response = requests.post(url = 'http://' + OF_Gateway_IP + ':' + OF_Gateway_Port + '/function/sha>
    #    "Subject": 'Negative Review Received',
    #    "Message": 'Review (ID = %i) of %s (ID = %i) received with negative results from sentiment a>
    #    event['reviewType'], int(event['productID']), int(event['customerID']), event['feedback'])
    #})

    payload = {
        'sentiment': event['sentiment'],
        'reviewType': event['reviewType'],
        'reviewID': event['reviewID'],
        'customerID': event['customerID'],
        'productID': event['productID'],
        'feedback': event['feedback']
    }
    response = Invoke_JSON(context, "SENTIMENT_DB", payload)
    return response[0], 200
