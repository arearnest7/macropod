from parliament import Context
from flask import Request
import json
import os
import datetime
import redis


def main(context: Context):
    print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "0" + "\n", flush=True)
    return b'a' * int(os.environ["LEN"]), 200
