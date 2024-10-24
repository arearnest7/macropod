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


def handle_handler(req):
    workflow_id = str(random.randint(0, 10000000))
    workflow_depth = 0
    workflow_width = 0
    if "workflow_id" in req:
        workflow_id = req["workflow_id"]
        workflow_depth = req["workflow_depth"]
        workflow_width = req["workflow_width"]
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
            fs.append(executor.submit(requests.post, url=os.environ["PIPELINED_CHECKSUM_PARTIAL"], json={"event": to_checksum, "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0}))
        if to_zip:
            fs.append(executor.submit(requests.post, url=os.environ["PIPELINED_ZIP_PARTIAL"], json={"event": to_zip, "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0}))
        if to_encrypt:
            fs.append(executor.submit(requests.post, url=os.environ["PIPELINED_ENCRYPT_PARTIAL"], json={"event": to_encrypt, "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0}))
    results = [f.result().text for f in fs]
    if to_checksum or to_zip:
        if to_checksum and "success" not in results[0]:
            to_checksum = []
        response = handle_handler({"manifest": new_manifest, "to_zip": to_checksum, "to_encrypt": to_zip, "workflow_id": workflow_id, "workflow_depth": workflow_depth, "workflow_width": 0})
        return response
    return "success"

def main(context: Context):
    if 'request' in context.keys():
        workflow_id = str(random.randint(0, 10000000))
        workflow_depth = 0
        workflow_width = 0
        if "workflow_id" in context.request.json:
            workflow_id = context.request.json["workflow_id"]
            workflow_depth = context.request.json["workflow_depth"]
            workflow_width = context.request.json["workflow_width"]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "0" + "\n", flush=True)
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
                fs.append(executor.submit(requests.post, url=os.environ["PIPELINED_CHECKSUM_PARTIAL"], json={"event": to_checksum, "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0}))
            if to_zip:
                fs.append(executor.submit(requests.post, url=os.environ["PIPELINED_ZIP_PARTIAL"], json={"event": to_zip, "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0}))
            if to_encrypt:
                fs.append(executor.submit(requests.post, url=os.environ["PIPELINED_ENCRYPT_PARTIAL"], json={"event": to_encrypt, "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0}))
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "1" + "\n", flush=True)
        results = [f.result().text for f in fs]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "2" + "\n", flush=True)
        if to_checksum or not to_zip:
            if not to_checksum and "success" not in results[0]:
                to_checksum = []
            response = handle_handler({"manifest": new_manifest, "to_zip": to_checksum, "to_encrypt": to_zip, "workflow_id": workflow_id, "workflow_depth": workflow_depth, "workflow_width": 0})
            print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "3" + "\n", flush=True)
            return response, 200
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "4" + "\n", flush=True)
        return "success", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
