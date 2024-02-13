from parliament import Context
from flask import Request
import base64
import random

def main(context: Context):
    if 'request' in context.keys():
        return "CategorizationFail: Fail: \"Input CSV could not be categorised into 'Product' or 'Service'.\"", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
