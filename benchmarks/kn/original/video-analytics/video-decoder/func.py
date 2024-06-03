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

def Recognise(workflow_id, workflow_depth, frame):
    result = requests.post(os.environ['VIDEO_RECOG'], json={"frame": base64.b64encode(frame).decode(), "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0}).text

    return result

def processFrames(videoBytes, workflow_id, workflow_depth):
    workflow_width = 0
    frames = decode(videoBytes)
    all_result_futures = []
    # send all requests
    frames = frames[0:6]
    ex = ThreadPoolExecutor(max_workers=6)
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "1" + "\n", flush=True)
    all_result_futures = ex.map(partial(Recognise, workflow_id, workflow_depth), frames)
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "2" + "\n", flush=True)
    results = ""
    for result in all_result_futures:
        results = results + result + ","

    return results

def Decode(request):
    videoBytes = b''
    videoBytes = base64.b64decode(request["video"].encode())
    results = processFrames(videoBytes, request["workflow_id"], request["workflow_depth"])
    return results

def main(context: Context):
    if 'request' in context.keys():
        workflow_id = str(random.randint(0, 10000000))
        workflow_depth = 0
        workflow_width = 0
        if "workflow_id" in context.request.json:
            workflow_id = context.request.json["workflow_id"]
            workflow_depth = context.request.json["workflow_depth"]
            workflow_width = context.request.json["workflow_width"]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "0" + "\n", flush=True)
        ret = Decode(context.request.json)
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "3" + "\n", flush=True)
        return ret, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
