from parliament import Context
from flask import Request
import json
import os
import requests

def main(context: Context):
    return requests.get(os.environ["DEST"], data=(b'a' * int(os.environ["LEN"]))).text, 200
