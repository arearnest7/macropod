import json
import os
from rpc import RPC

def function_handler(context):
    if context["request_type"] != "GRPC":
        return str(RPC(os.environ["TEST"], ["TEST"], context["workflow_id"])), 200
    return context["request"], 200
