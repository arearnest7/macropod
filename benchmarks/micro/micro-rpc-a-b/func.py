from parliament import Context
from flask import Request
import json
import os
import datetime
import redis
import requests


def b_main(a):
    return b'a' * int(os.environ["LEN_B"])

def main(context: Context):
    print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "0" + "\n", flush=True)
    ret = b_main(b'a' * int(os.environ["LEN_A"]))
    print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "1" + "\n", flush=True)
    return ret, 200
