from parliament import Context
from flask import Request
import base64
import requests
import json
import pprint
import csv
import os
import sys
from pymongo import MongoClient
from urllib.parse import quote_plus
import redis
import random

pp = pprint.PrettyPrinter(indent=4)

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

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


def sfail_handler(req):
    return "SentimentFail: Fail: \"Sentiment Analysis Failed!\""

def sns_handler(req):
    event = req
    #response = requests.get(url = 'http://' + OF_Gateway_IP + ':' + OF_Gateway_Port + '/function/shasum', json={
    #    "Subject": 'Negative Review Received',
    #    "Message": 'Review (ID = %i) of %s (ID = %i) received with negative results from sentiment analysis. Feedback from Customer (ID = %i): "%s"' % (int(event['reviewID']),
    #    event['reviewType'], int(event['productID']), int(event['customerID']), event['feedback'])
    #})

    return db_handler({
        'sentiment': event['sentiment'],
        'reviewType': event['reviewType'],
        'reviewID': event['reviewID'],
        'customerID': event['customerID'],
        'productID': event['productID'],
        'feedback': event['feedback'],
    })

def service_result_handler(req):
    event = req
    results = ""
    if event["sentiment"] == "POSITIVE" or event["sentiment"] == "NEUTRAL":
        results = db_handler(event)
    elif event["sentiment"] == "NEGATIVE":
        results = sns_handler(event)
    else:
        results = sfail_handler(event)

    return results

def product_result_handler(req):
    event = req
    results = ""
    if event["sentiment"] == "POSITIVE" or event["sentiment"] == "NEUTRAL":
        results = db_handler(event)
    elif event["sentiment"] == "NEGATIVE":
        results = sns_handler(event)
    else:
        results = sfail_handler(event)

    return results

def cfail_handler(req):
    return "CategorizationFail: Fail: \"Input CSV could not be categorised into 'Product' or 'Service'.\""


def product_sentiment_handler(req):
    event = req
    feedback = event['feedback']
    response = {"polarity": -0.66}
    if response['polarity'] > 0.5:
        sentiment = "POSITIVE"
    elif response['polarity'] < -0.5:
        sentiment = "NEGATIVE"
    else:
        sentiment = "NEUTRAL"

    return product_result_handler({
        'sentiment': sentiment,
        'reviewType': event['reviewType'],
        'reviewID': event['reviewID'],
        'customerID': event['customerID'],
        'productID': event['productID'],
        'feedback': event['feedback'],
    })

def service_sentiment_handler(req):
    event = req
    feedback = event['feedback']
    response = {"polarity": -0.66}
    if response['polarity'] > 0.5:
        sentiment = "POSITIVE"
    elif response['polarity'] < -0.5:
        sentiment = "NEGATIVE"
    else:
        sentiment = "NEUTRAL"

    return service_result_handler({
        'sentiment': sentiment,
        'reviewType': event['reviewType'],
        'reviewID': event['reviewID'],
        'customerID': event['customerID'],
        'productID': event['productID'],
        'feedback': event['feedback'],
    })

def product_or_service_handler(req):
    event = req
    results = ""
    if event["reviewType"] == "Product":
        results = product_sentiment_handler(event)
    elif event["reviewType"] == "Service":
        results = service_sentiment_handler(event)
    else:
        results = cfail_handler(event)
    return results

def read_csv_handler(req):
    event = req

    bucket_name = event['bucket_name']
    file_key = event['file_key']
    response = open(file_key, 'r').read()

    lines = response.split('\n')

    for row in csv.DictReader(lines):
        return product_or_service_handler(row)

def main(context: Context):
    if 'request' in context.keys():
        event = context.request.json

        bucket_name = event['Records'][0]['s3']['bucket']['name']
        file_key = event['Records'][0]['s3']['object']['key']

        input= {
                'bucket_name': bucket_name,
                'file_key': file_key
            }
        response = read_csv_handler(input)
        return response, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
