from parliament import Context
from flask import Request
import base64
import requests
import json
import os
import sys
from pymongo import MongoClient
from urllib.parse import quote_plus
import random
import datetime
import redis


def db_handler(req):
    event = req
    #client = MongoClient(host=os.environ["MONGO_HOST"])
    #db = client['sentiment']
    #table = ""
    #if event['reviewType'] == 'Product':
    #    table = db.products
    #elif event['reviewType'] == 'Service':
    #    table = db.services
    #else:
    #    raise Exception("Input review is neither Product nor Service")
    #Item = {
    #   'reviewID': event['reviewID'],
    #    'customerID': event['customerID'],
    #    'productID': event['productID'],
    #    'feedback': event['feedback'],
    #    'sentiment': event['sentiment']
    #}
    #response = {"response": str(table.insert_one(Item).inserted_id)}

    response = {}
    response['reviewType'] = event['reviewType']
    response['reviewID'] = event['reviewID']
    response['customerID'] = event['customerID']
    response['productID'] = event['productID']
    response['feedback'] = event['feedback']
    response['sentiment'] = event['sentiment']

    return json.dumps(response)

def sns_handler(req):
    event = req
    #response = requests.get(url = 'http://' + OF_Gateway_IP + ':' + OF_Gateway_Port + '/function/sha>    #    "Subject": 'Negative Review Received',
    #    "Message": 'Review (ID = %i) of %s (ID = %i) received with negative results from sentiment a>    #    event['reviewType'], int(event['productID']), int(event['customerID']), event['feedback'])
    #})

    return db_handler({
        'sentiment': event['sentiment'],
        'reviewType': event['reviewType'],
        'reviewID': event['reviewID'],
        'customerID': event['customerID'],
        'productID': event['productID'],
        'feedback': event['feedback'],
    })

def main(context: Context):
    if 'request' in context.keys():
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "0" + "\n", flush=True)
        event = context.request.json
        results = ""
        if event["sentiment"] == "POSITIVE" or event["sentiment"] == "NEUTRAL":
            results = db_handler(event)
        elif event["sentiment"] == "NEGATIVE":
            results = sns_handler(event)
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "1" + "\n", flush=True)
        return results, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
