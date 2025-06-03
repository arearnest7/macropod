from parliament import Context
from flask import Request
import base64
import requests
import os
import json
import string
from concurrent.futures import ThreadPoolExecutor
import random
import datetime
import redis


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
                fs.append(executor.submit(requests.post, url=os.environ["PIPELINED_CHECKSUM"], json={"event": to_checksum, "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0}))
            if to_zip:
                fs.append(executor.submit(requests.post, url=os.environ["PIPELINED_ZIP"], json={"event": to_zip, "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0}))
            if to_encrypt:
                fs.append(executor.submit(requests.post, url=os.environ["PIPELINED_ENCRYPT"], json={"event": to_encrypt, "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0}))
        results = [f.result().text for f in fs]
        if to_checksum or to_zip:
            if to_checksum and "success" not in results[0]:
                to_checksum = []
            response = requests.post(url=os.environ["PIPELINED_MAIN"], json={"manifest": new_manifest, "to_zip": to_checksum, "to_encrypt": to_zip, "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0})
            return response.text, 200
        return "success", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
