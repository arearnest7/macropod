from parliament import Context
from flask import Request
import json
import os
import datetime
import redis
import requests

if "LOGGING_NAME" in os.environ:
    loggingClient = redis.Redis(host=os.environ['LOGGING_IP'], password=os.environ['LOGGING_PASSWORD'])

def b_main(a):
    return b'a' * int(os.environ["LEN_B"])

def main(context: Context):
    if "LOGGING_NAME" in os.environ:
        loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n")
    ret = b_main(b'a' * int(os.environ["LEN_A"]))
    if "LOGGING_NAME" in os.environ:
        loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n")
    return ret, 200
