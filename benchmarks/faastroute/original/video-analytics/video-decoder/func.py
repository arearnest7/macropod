from rpc import RPC
from __future__ import print_function
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

def Recognise(frame, context):
    result = RPC(context, os.environ['VIDEO_RECOG'], [frame])[0]

    return result

def processFrames(videoBytes):
    frames = decode(videoBytes)
    all_result_futures = []
    # send all requests
    frames = frames[0:6]
    ex = ThreadPoolExecutor(max_workers=6)
    all_result_futures = ex.map(Recognise, frames, context)
    results = ""
    for result in all_result_futures:
        results = results + result + ","

    return results

def Decode(request):
    videoBytes = b''
    videoBytes = request
    results = processFrames(videoBytes)
    return results

def function_handler(context):
    if context["InvokeType"] == "GRPC":
        ret = Decode(context["Request"])
        return ret, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
