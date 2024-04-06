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
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "10" + "\n", flush=True)
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
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "11" + "\n", flush=True)
        results = [f.result().text for f in fs]
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "12" + "\n", flush=True)
        if to_checksum or to_zip:
            if to_checksum and "success" not in results[0]:
                to_checksum = []
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "13" + "\n", flush=True)
            response = requests.get(url=os.environ["PIPELINED_MAIN"], json={"manifest": new_manifest, "to_zip": to_checksum, "to_encrypt": to_zip})
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "14" + "\n", flush=True)
            return response.text, 200
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "15" + "\n", flush=True)
        return "success", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
