import base64
import random

def function_handler(context):
    if context["is_json"]:
        return "CategorizationFail: Fail: \"Input CSV could not be categorised into 'Product' or 'Service'.\"", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200