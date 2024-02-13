from parliament import Context
from flask import Request
import base64
import random

def main(context: Context):
    if 'request' in context.keys():
        return "SentimentFail: Fail: \"Sentiment Analysis Failed!\"", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
