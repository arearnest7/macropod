from parliament import Context
from flask import Request
import json
import os
import datetime
import redis
import requests
import datetime

def main(context: Context):
    print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "0" + "\n", flush=True)
    ret = requests.get(os.environ["DEST"], data=(b'a' * int(os.environ["LEN"]))).text
    print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "1" + "\n", flush=True)
    return ret, 200
