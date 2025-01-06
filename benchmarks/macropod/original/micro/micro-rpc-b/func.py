import json
import os
from rpc import RPC

def FunctionHandler(context):
    return str(b'a' * int(os.environ["LEN"])), 200
