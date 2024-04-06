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


def main(context: Context):
    if 'request' in context.keys():
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "0" + "\n", flush=True)
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

        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "1" + "\n", flush=True)
        return json.dumps(response), 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
