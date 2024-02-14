from rpc import RPC
import base64
import random

def function_handler(context):
    if context["request_type"] == "GRPC":
        return "SentimentFail: Fail: \"Sentiment Analysis Failed!\"", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
