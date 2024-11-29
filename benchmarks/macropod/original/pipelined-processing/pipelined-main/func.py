from rpc import RPC
import base64
import requests
import os
import json
import string
from concurrent.futures import ThreadPoolExecutor
import random

def FunctionHandler(context):
    body = json.loads(context["Request"])
    if body["manifest"]:
        to_checksum = body["manifest"][0]
    else:
        to_checksum = []
    to_zip = body["to_zip"]
    to_encrypt = body["to_encrypt"]
    if len(body["manifest"]) > 1:
        new_manifest = body["manifest"][1:]
    else:
        new_manifest = []

    fs = []
    with ThreadPoolExecutor(max_workers=3) as executor:
        if to_checksum:
            fs.append(executor.submit(RPC, context, os.environ["PIPELINED_CHECKSUM"], [json.dumps(to_checksum).encode()]))
        if to_zip:
            fs.append(executor.submit(RPC, context, os.environ["PIPELINED_ZIP"], [json.dumps(to_zip).encode()]))
        if to_encrypt:
            fs.append(executor.submit(RPC, context, os.environ["PIPELINED_ENCRYPT"], [json.dumps(to_encrypt).encode()]))
    results = [f.result()[0] for f in fs]
    if to_checksum or to_zip:
        if to_checksum and "success" not in results[0]:
            to_checksum = []
        response = RPC(context, os.environ["PIPELINED_MAIN"], [json.dumps({"manifest": new_manifest, "to_zip": to_checksum, "to_encrypt": to_zip}).encode()])[0]
        return response, 200
    return "success", 200
