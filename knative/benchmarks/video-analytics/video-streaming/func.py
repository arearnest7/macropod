from parliament import Context
from flask import Request
import json
import os
import requests
import base64
import datetime
import random

def main(context: Context):
    videoFile = open("reference/" + context.request.json["video"], "rb")
    videoFragment = videoFile.read()
    videoFile.close()
    headers = {'Content-Type': 'application/octet-stream'}
    ret = requests.post(os.environ['VIDEO_DECODER'], data=videoFragment, headers=headers).text
    return ret, 200
