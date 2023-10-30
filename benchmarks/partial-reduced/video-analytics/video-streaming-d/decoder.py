from parliament import Context
from flask import Request
import json
from __future__ import print_function

import pickle
import sys
import os
# adding python tracing and storage sources to the system path
sys.path.insert(0, os.getcwd() + '/../proto/')
sys.path.insert(0, os.getcwd() + '/../../../../utils/tracing/python')
sys.path.insert(0, os.getcwd() + '/../../../../utils/storage/python')
import tracing
import videoservice_pb2_grpc
import videoservice_pb2
import destination as XDTdst
import source as XDTsrc
import utils as XDTutil

import cv2
import grpc
import tempfile
import argparse
import boto3
import logging as log
import socket
import requests

from concurrent import futures

with open('/etc/secret-volume/video-recog-partial', 'r') as f:
    video-recog = f.read()

def decode(bytes):
    temp = tempfile.NamedTemporaryFile(suffix=".mp4")
    temp.write(bytes)
    temp.seek(0)

    all_frames = []
    with tracing.Span("Decode frames"):
        vidcap = cv2.VideoCapture(temp.name)
        for i in range(int(os.getenv('DecoderFrames', int(args.frames)))):
            success,image = vidcap.read()
            all_frames.append(cv2.imencode('.jpg', image)[1].tobytes())

    return all_frames

def Recognise(frame):
    result = requests.get(url = video-recog + ":80", data = json.dumps(frame)).text

    return result

def processFrames(videoBytes):
    frames = decode(videoBytes)
    with tracing.Span("Recognise all frames"):
        all_result_futures = []
        # send all requests
        decoderFrames = int(os.getenv('DecoderFrames', 6))
        frames = frames[0:decoderFrames]
        if os.getenv('CONCURRENT_RECOG', "false").lower() == "false":
            # concat all results
            for frame in frames:
                all_result_futures.append(Recognise(frame))
        else:
            ex = futures.ThreadPoolExecutor(max_workers=decoderFrames)
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

def main():
    ret = Decode(json.loads(sys.argv[1]))
    return ret, 200
main()
