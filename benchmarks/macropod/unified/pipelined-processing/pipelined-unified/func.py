from rpc import Invoke
import base64
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

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])


def checksum_handler(req):
    event = req
    data = open("original-" + event[0], 'rb').read()
    md5 = hashlib.md5(data).hexdigest()
    if event[1] == md5:
        #redisClient.set("checksumed-" + event[0], data)
        return "success"
    return "failed"

def zip_handler(req):
    event = req
    data = open("checksumed-" + event, 'rb').read()
    with open("/tmp/" + event, "wb") as f:
        f.write(data)
    with ZipFile('/tmp/zip.zip', 'w') as zip:
        zip.write("/tmp/" + event)
    zip.close()
    with open("/tmp/zip.zip", "rb") as f:
        data = f.read()
    #redisClient.set("ziped-" + event, data)
    return "success"

def encrypt_handler(req):
    event = req
    data = open("ziped-" + event, 'rb').read()
    with open("/tmp/" + event + ".zip", "wb") as f:
        f.write(data)
    key = Fernet.generate_key()
    with open('/tmp/key.key', 'wb') as filekey:
        filekey.write(key)
    filekey.close()
    fernet = Fernet(key)
    data = ""
    with open("/tmp/" + event + ".zip", "rb") as file:
        data = file.read()
    file.close()
    encrypted_data = fernet.encrypt(data)
    #redisClient.set("encrypted-" + event, encrypted_data)
    #redisClient.set("encrypted-key-" + event, key)
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
            fs.append(executor.submit(checksum_handler, to_checksum))
        if to_zip:
            fs.append(executor.submit(zip_handler, to_zip[0]))
        if to_encrypt:
            fs.append(executor.submit(encrypt_handler, to_encrypt[0]))
    results = [f.result() for f in fs]
    if to_checksum or to_zip:
        if to_checksum and "success" not in results[0]:
            to_checksum = []
        return handle_handler({"manifest": new_manifest, "to_zip": to_checksum, "to_encrypt": to_zip})
    return "success"

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
            fs.append(executor.submit(checksum_handler, to_checksum))
        if to_zip:
            fs.append(executor.submit(zip_handler, to_zip[0]))
        if to_encrypt:
            fs.append(executor.submit(encrypt_handler, to_encrypt[0]))
    results = [f.result() for f in fs]
    if to_checksum or to_zip:
        if to_checksum and "success" not in results[0]:
            to_checksum = []
        response = handle_handler({"manifest": new_manifest, "to_zip": to_checksum, "to_encrypt": to_zip})
        return response, 200
    return "success", 200
