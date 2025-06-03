from rpc import Invoke
import base64
import os
import json
import string
from cryptography.fernet import Fernet
import redis
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def FunctionHandler(context):
    event = context["TEXT"]
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
    return "success", 200
