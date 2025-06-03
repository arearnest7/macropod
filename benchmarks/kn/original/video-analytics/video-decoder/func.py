from __future__ import print_function
from parliament import Context
from flask import Request
import json

import pickle
import sys
import os

import cv2
import tempfile
import argparse
import socket
import requests
import base64

from concurrent.futures import ThreadPoolExecutor
from functools import partial
import datetime
import redis
import random

def decode(bytes):
    temp = tempfile.NamedTemporaryFile(suffix=".mp4")
    temp.write(bytes)
    temp.seek(0)

    all_frames = []
    vidcap = cv2.VideoCapture(temp.name)
    for i in range(int(6)):
        success,image = vidcap.read()
        all_frames.append(cv2.imencode('.jpg', image)[1].tobytes())

    return all_frames

def Recognise(frame):
    result = requests.post(os.environ['VIDEO_RECOG'], json={"frame": base64.b64encode(frame).decode()}).text

    return result

def processFrames(videoBytes):
    workflow_width = 0
    frames = decode(videoBytes)
    all_result_futures = []
    # send all requests
    frames = frames[0:6]
    ex = ThreadPoolExecutor(max_workers=6)
    all_result_futures = ex.map(partial(Recognise, frames)
    results = ""
    for result in all_result_futures:
        results = results + result + ","

    return results

def Decode(request):
    videoBytes = b''
    videoBytes = base64.b64decode(request["video"].encode())
    results = processFrames(videoBytes)
    return results

def main(context: Context):
    if 'request' in context.keys():
        ret = Decode(context.request.json)
        return ret, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
