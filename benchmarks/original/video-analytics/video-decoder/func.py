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
import logging as log
import socket
import requests

sys.path.insert(0, os.getcwd() + 'utils/tracing/python')
import tracing

from concurrent.futures import ThreadPoolExecutor

def decode(bytes):
    temp = tempfile.NamedTemporaryFile(suffix=".mp4")
    temp.write(bytes)
    temp.seek(0)

    all_frames = []
    with tracing.Span("Decode frames"):
        vidcap = cv2.VideoCapture(temp.name)
        for i in range(int(6)):
            success,image = vidcap.read()
            all_frames.append(cv2.imencode('.jpg', image)[1].tobytes())

    return all_frames

def Recognise(frame):
    result = requests.get(url = os.environ['VIDEO_RECOG'] + ":8080", data = json.dumps(frame)).text

    return result

def processFrames(videoBytes):
    frames = decode(videoBytes)
    with tracing.Span("Recognise all frames"):
        all_result_futures = []
        # send all requests
        frames = frames[0:6]
        if os.getenv('CONCURRENT_RECOG', "false").lower() == "false":
            # concat all results
            for frame in frames:
                all_result_futures.append(Recognise(frame))
        else:
            ex = ThreadPoolExecutor(max_workers=6)
            all_result_futures = ex.map(Recognise, frames)
        log.info("returning result of frame classification")
        results = ""
        for result in all_result_futures:
            results = results + result + ","

        return results

def Decode(request):
    log.info("Decoder recieved a request")

    videoBytes = b''
    log.info("Inline video decode. Decoding frames.")
    videoBytes = request.video
    results = processFrames(videoBytes)
    return results

def main(context: Context):
    if 'request' in context.keys():
        ret = Decode(context.request.json)
        return ret, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
