import base64
import random

def function_handler(context):
    if context["is_json"]:
        return "SentimentFail: Fail: \"Sentiment Analysis Failed!\"", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
