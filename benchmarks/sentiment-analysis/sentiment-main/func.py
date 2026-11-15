from rpc import Invoke_JSON
import base64
import json
import pprint
import random
import os

pp = pprint.PrettyPrinter(indent=4)

def FunctionHandler(context):
    event = context["JSON"]

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
    response = Invoke_JSON(context, "SENTIMENT_READ_CSV", input)
    return response[0], 200
