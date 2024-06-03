import base64
import requests
import json
import sys
import os
from pymongo import MongoClient
from urllib.parse import quote_plus
import random
import datetime

def function_handler(context):
    if context["is_json"]:
        workflow_id = context["request"]["workflow_id"]
        workflow_depth = context["request"]["workflow_depth"]
        workflow_width = context["request"]["workflow_width"]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "3" + "\n", flush=True)
        event = context["request"]

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

        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "4" + "\n", flush=True)
        return json.dumps(response), 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
