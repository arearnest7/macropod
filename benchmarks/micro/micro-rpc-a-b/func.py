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
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "0" + "\n", flush=True)
    ret = b_main(b'a' * int(os.environ["LEN_A"]))
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "1" + "\n", flush=True)
    return ret, 200
