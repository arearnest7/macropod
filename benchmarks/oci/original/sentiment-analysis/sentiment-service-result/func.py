import base64
import requests
import json
import random
import os
import datetime

def function_handler(context):
    if context["is_json"]:
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "10" + "\n", flush=True)
        event = context["request"]
        results = ""
        if event["sentiment"] == "POSITIVE" or event["sentiment"] == "NEUTRAL":
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "11" + "\n", flush=True)
            results = requests.get(url=os.environ["SENTIMENT_DB"], json=event)
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "12" + "\n", flush=True)
        elif event["sentiment"] == "NEGATIVE":
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "13" + "\n", flush=True)
            results = requests.get(url=os.environ["SENTIMENT_SNS"], json=event)
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "14" + "\n", flush=True)
        else:
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "15" + "\n", flush=True)
            results = requests.get(url=os.environ["SENTIMENT_SFAIL"], json=event)
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "16" + "\n", flush=True)

        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "17" + "\n", flush=True)
        return results.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
