from rpc import Invoke
from rpc import Invoke_JSON
import base64
import requests
import os
import json
import string
from concurrent.futures import ThreadPoolExecutor
import random

def FunctionHandler(context):
    body = context["JSON"]
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
            fs.append(executor.submit(Invoke_JSON, context, "PIPELINED_CHECKSUM", to_checksum))
        if to_zip:
            fs.append(executor.submit(Invoke, context, "PIPELINED_ZIP", to_zip[0]))
        if to_encrypt:
            fs.append(executor.submit(Invoke, context, "PIPELINED_ENCRYPT", to_encrypt[0]))
    results = [f.result()[0] for f in fs]
    if to_checksum or to_zip:
        if to_checksum and "success" not in results[0]:
            to_checksum = []
        response, code = Invoke(context, "PIPELINED_MAIN", {"manifest": new_manifest, "to_zip": to_checksum, "to_encrypt": to_zip})
        return response, 200
    return "success", 200
