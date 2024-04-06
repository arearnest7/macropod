from parliament import Context
from flask import Request
import base64
import requests
import os
import json
import string
from concurrent.futures import ThreadPoolExecutor
import hashlib
from zipfile import ZipFile
from cryptography.fernet import Fernet
import datetime
import redis
import random

redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])


def checksum_handler(req):
    event = req["event"]
    data = redisClient.get("original-" + event[0])
    md5 = hashlib.md5(data).hexdigest()
    if event[1] == md5:
        redisClient.set("checksumed-" + event[0], data)
        return "success"
    return "failed"

def zip_handler(req):
    event = req["event"]
    data = redisClient.get("checksumed-" + event[0])
    with open("/tmp/" + event[0], "wb") as f:
        f.write(data)
    with ZipFile('/tmp/zip.zip', 'w') as zip:
        zip.write("/tmp/" + event[0])
    zip.close()
    with open("/tmp/zip.zip", "rb") as f:
        data = f.read()
    redisClient.set("ziped-" + event[0], data)
    return "success"

def encrypt_handler(req):
    event = req["event"]
    data = redisClient.get("ziped-" + event[0])
    with open("/tmp/" + event[0] + ".zip", "wb") as f:
        f.write(data)
    key = Fernet.generate_key()
    with open('/tmp/key.key', 'wb') as filekey:
        filekey.write(key)
    filekey.close()
    fernet = Fernet(key)
    data = ""
    with open("/tmp/" + event[0] + ".zip", "rb") as file:
        data = file.read()
    file.close()
    encrypted_data = fernet.encrypt(data)
    redisClient.set("encrypted-" + event[0], encrypted_data)
    redisClient.set("encrypted-key-" + event[0], key)
    return "success"

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
            fs.append(executor.submit(checksum_handler, {"event": to_checksum}))
        if to_zip:
            fs.append(executor.submit(zip_handler, {"event": to_zip}))
        if to_encrypt:
            fs.append(executor.submit(encrypt_handler, {"event": to_encrypt}))
    results = [f.result() for f in fs]
    if to_checksum or to_zip:
        if to_checksum and "success" not in results[0]:
            to_checksum = []
        return handle_handler({"manifest": new_manifest, "to_zip": to_checksum, "to_encrypt": to_zip})
    return "success"

def main(context: Context):
    if 'request' in context.keys():
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "0" + "\n", flush=True)
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
                fs.append(executor.submit(checksum_handler, {"event": to_checksum}))
            if to_zip:
                fs.append(executor.submit(zip_handler, {"event": to_zip}))
            if to_encrypt:
                fs.append(executor.submit(encrypt_handler, {"event": to_encrypt}))
        results = [f.result() for f in fs]
        if to_checksum or to_zip:
            if to_checksum and "success" not in results[0]:
                to_checksum = []
            print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "1" + "\n", flush=True)
            return handle_handler({"manifest": new_manifest, "to_zip": to_checksum, "to_encrypt": to_zip})
        print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "2" + "\n", flush=True)
        return "success", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
