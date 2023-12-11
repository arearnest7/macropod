from parliament import Context
from flask import Request
import base64
import requests
import os
import json
import string
from concurrent.futures import ThreadPoolExecutor
import random

def handle_handler(req):
    if req["manifest"]:
        to_checksum = req["manifest"][0]
    else:
        to_checksum = []
    to_zip = req["to_zip"]
    to_encrypt = req["to_encrypt"]
    if len(req["manifest"]) > 1:
        new_manifest = req["manifest"][1:]
    else:
        new_manifest = []

    fs = []
    with ThreadPoolExecutor(max_workers=3) as executor:
        if to_checksum:
            fs.append(executor.submit(requests.get, url=os.environ["PIPELINED_CHECKSUM_PARTIAL"], json={"event": to_checksum}))
        if to_zip:
            fs.append(executor.submit(requests.get, url=os.environ["PIPELINED_ZIP_PARTIAL"], json={"event": to_zip}))
        if to_encrypt:
            fs.append(executor.submit(requests.get, url=os.environ["PIPELINED_ENCRYPT_PARTIAL"], json={"event": to_encrypt}))
    results = [f.result().text for f in fs]
    if to_checksum or to_zip:
        if to_checksum and "success" not in results[0]:
            to_checksum = []
        response = handle_handler({"manifest": new_manifest, "to_zip": to_checksum, "to_encrypt": to_zip})
        return response
    return "success"

def main(context: Context):
    if 'request' in context.keys():
        if context.request.json["manifest"]:
            to_checksum = context.request.json["manifest"][0]
        else:
            to_checksum = []
        to_zip = context.request.json["to_zip"]
        to_encrypt = context.request.json["to_encrypt"]
        if len(context.request.json["manifest"]) > 1:
            new_manifest = context.request.json["manifest"][1:]
        else:
            new_manifest = []

        fs = []
        with ThreadPoolExecutor(max_workers=3) as executor:
            if to_checksum:
                fs.append(executor.submit(requests.get, url=os.environ["PIPELINED_CHECKSUM_PARTIAL"], json={"event": to_checksum}))
            if to_zip:
                fs.append(executor.submit(requests.get, url=os.environ["PIPELINED_ZIP_PARTIAL"], json={"event": to_zip}))
            if to_encrypt:
                fs.append(executor.submit(requests.get, url=os.environ["PIPELINED_ENCRYPT_PARTIAL"], json={"event": to_encrypt}))
        results = [f.result().text for f in fs]
        if to_checksum or not to_zip:
            if not to_checksum and "success" not in results[0]:
                to_checksum = []
            response = handle_handler({"manifest": new_manifest, "to_zip": to_checksum, "to_encrypt": to_zip})
            return response, 200
        return "success", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
