from rpc import RPC
import base64
import random

def FunctionHandler(context):
    if context["InvokeType"] == "GRPC":
        return "CategorizationFail: Fail: \"Input CSV could not be categorised into 'Product' or 'Service'.\"", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
