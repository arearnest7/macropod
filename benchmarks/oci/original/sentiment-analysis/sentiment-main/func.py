import base64
import requests
import json
import pprint
import random
import os
import datetime

pp = pprint.PrettyPrinter(indent=4)

def function_handler(context):
    if context["is_json"]:
        workflow_id = context["request"]["workflow_id"]
        workflow_depth = context["request"]["workflow_depth"]
        workflow_width = context["request"]["workflow_width"]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "3" + "\n", flush=True)
        event = context["request"]

        try:
            pp
        except NameError:
            pp = pprint.PrettyPrinter(indent=4)

        bucket_name = event['Records'][0]['s3']['bucket']['name']
        file_key = event['Records'][0]['s3']['object']['key']

        input= {
                'bucket_name': bucket_name,
                'file_key': file_key,
                "workflow_id": workflow_id,
                "workflow_depth": workflow_depth + 1,
                "workflow_width": 0
            }
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "4" + "\n", flush=True)
        response = requests.post(url=os.environ["SENTIMENT_READ_CSV"], json=input)
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "5" + "\n", flush=True)

        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
