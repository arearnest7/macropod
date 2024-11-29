from rpc import RPC
import base64
import requests
import json
import sys
import os
from pymongo import MongoClient
from urllib.parse import quote_plus
import random

def FunctionHandler(context):
    event = json.loads(context["Request"])

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

    return json.dumps(response), 200
