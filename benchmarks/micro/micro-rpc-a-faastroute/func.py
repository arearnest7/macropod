import json
import os
from rpc import RPC

def FunctionHandler(context):
    return str(RPC(context, os.environ["DEST"], [(b'a' * int(os.environ["LEN"]))])), 200
