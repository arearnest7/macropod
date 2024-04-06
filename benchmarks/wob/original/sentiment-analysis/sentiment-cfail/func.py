from parliament import Context
from flask import Request
import base64
import random
import os
import datetime
import redis


def main(context: Context):
    if 'request' in context.keys():
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "0" + "\n", flush=True)
        return "CategorizationFail: Fail: \"Input CSV could not be categorised into 'Product' or 'Service'.\"", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
