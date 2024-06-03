import base64
import requests
import os
import json
import string
from concurrent.futures import ThreadPoolExecutor
import random
import datetime

def function_handler(context):
    if context["is_json"]:
        workflow_id = context["request"]["workflow_id"]
        workflow_depth = context["request"]["workflow_depth"]
        workflow_width = context["request"]["workflow_width"]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "3" + "\n", flush=True)
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
                fs.append(executor.submit(requests.post, url=os.environ["PIPELINED_CHECKSUM"], json={"event": to_checksum, "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0}))
            if to_zip:
                fs.append(executor.submit(requests.post, url=os.environ["PIPELINED_ZIP"], json={"event": to_zip, "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0}))
            if to_encrypt:
                fs.append(executor.submit(requests.post, url=os.environ["PIPELINED_ENCRYPT"], json={"event": to_encrypt, "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0}))
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "4" + "\n", flush=True)
        results = [f.result().text for f in fs]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "5" + "\n", flush=True)
        if to_checksum or to_zip:
            if to_checksum and "success" not in results[0]:
                to_checksum = []
            print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "6" + "\n", flush=True)
            response = requests.post(url=os.environ["PIPELINED_MAIN"], json={"manifest": new_manifest, "to_zip": to_checksum, "to_encrypt": to_zip, "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0})
            print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "7" + "\n", flush=True)
            return response.text, 200
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "8" + "\n", flush=True)
        return "success", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
