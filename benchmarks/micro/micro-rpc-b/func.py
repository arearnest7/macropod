from parliament import Context
from flask import Request
import json
import os
import datetime
import redis


def main(context: Context):
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + "0" + "," + "0" + "," + "0" + "," + "HTTP" + "," + "0" + "\n", flush=True)
    return b'a' * int(os.environ["LEN"]), 200
