from parliament import Context
from flask import Request
import base64
import requests
import json
import sys
import os
from pymongo import MongoClient
from urllib.parse import quote_plus
import random
import datetime
import redis

if "LOGGING_NAME" in os.environ:
    loggingClient = redis.Redis(host=os.environ['LOGGING_IP'], password=os.environ['LOGGING_PASSWORD'])

def main(context: Context):
    if 'request' in context.keys():
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n")
        event = context.request.json

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

        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n")

        return json.dumps(response), 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
