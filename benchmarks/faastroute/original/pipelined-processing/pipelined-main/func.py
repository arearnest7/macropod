from rpc import RPC
import base64
import requests
import os
import json
import string
from concurrent.futures import ThreadPoolExecutor
import random

def function_handler(context):
    if context["is_json"]:
        if context["request"]["manifest"]:
            to_checksum = context["request"]["manifest"][0]
        else:
            to_checksum = []
        to_zip = context["request"]["to_zip"]
        to_encrypt = context["request"]["to_encrypt"]
        if len(context["request"]["manifest"]) > 1:
            new_manifest = context["request"]["manifest"][1:]
        else:
            new_manifest = []

        fs = []
        with ThreadPoolExecutor(max_workers=3) as executor:
            if to_checksum:
                fs.append(executor.submit(requests.get, url=os.environ["PIPELINED_CHECKSUM"], json={"event": to_checksum}))
            if to_zip:
                fs.append(executor.submit(requests.get, url=os.environ["PIPELINED_ZIP"], json={"event": to_zip}))
            if to_encrypt:
                fs.append(executor.submit(requests.get, url=os.environ["PIPELINED_ENCRYPT"], json={"event": to_encrypt}))
        results = [f.result().text for f in fs]
        if to_checksum or to_zip:
            if to_checksum and "success" not in results[0]:
                to_checksum = []
            response = requests.get(url=os.environ["PIPELINED_MAIN"], json={"manifest": new_manifest, "to_zip": to_checksum, "to_encrypt": to_zip})
            return response.text, 200
        return "success", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
