import json
import os
from rpc import RPC

def FunctionHandler(context):
    if "TEST" in os.environ:
        return str(RPC(context, os.environ["TEST"], [(b'A' * 10000000)])), 200
    return str(context["Request"]), 200
