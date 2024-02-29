from rpc import RPC
import base64
import requests
import json
import pprint
import random
import os

pp = pprint.PrettyPrinter(indent=4)

def function_handler(context):
    if context["request_type"] != "GRPC":
        event = context["request"]

        try:
            pp
        except NameError:
            pp = pprint.PrettyPrinter(indent=4)

        bucket_name = event['Records'][0]['s3']['bucket']['name']
        file_key = event['Records'][0]['s3']['object']['key']

        input= {
                'bucket_name': bucket_name,
                'file_key': file_key
            }
        response = RPC(context, os.environ["SENTIMENT_READ_CSV"], [json.dumps(input)])[0]

        return response, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
