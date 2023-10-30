from parliament import Context
from flask import Request
import json
from __future__ import print_function
from torchvision import transforms
from PIL import Image
import torch
import torchvision.models as models

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
import io

from concurrent import futures

# Load model
model = models.squeezenet1_1(pretrained=True)

labels_fd = open('imagenet_labels.txt', 'r')
labels = []
for i in labels_fd:
    labels.append(i)
labels_fd.close()

def preprocessImage(imageBytes):
    with tracing.Span("preprocess"):
        img = Image.open(io.BytesIO(imageBytes))

        transform = transforms.Compose([
            transforms.Resize(256),
            transforms.CenterCrop(224),
            transforms.ToTensor(),
            transforms.Normalize(
                mean=[0.485, 0.456, 0.406],
                std=[0.229, 0.224, 0.225]
            )
        ])

        img_t = transform(img)
        return torch.unsqueeze(img_t, 0)


def infer(batch_t):
    with tracing.Span("infer"):
        # Set up model to do evaluation
        model.eval()

        # Run inference
        with torch.no_grad():
            out = model(batch_t)

        # Print top 5 for logging
        _, indices = torch.sort(out, descending=True)
        percentages = torch.nn.functional.softmax(out, dim=1)[0] * 100
        for idx in indices[0][:5]:
            log.info('\tLabel: %s, percentage: %.2f' % (labels[idx], percentages[idx].item()))

        # make comma-seperated output of top 100 label
        out = ""
        for idx in indices[0][:100]:
            out = out + labels[idx] + ","
        return out

def Recognise_2(request):
    log.info("received a call")

    # get the frame from s3 or inline
    frame = None
    frame = request.frame

    log.info("performing image recognition on frame")
    classification = infer(preprocessImage(frame))
    log.info("object recogintion successful")
    return classification

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
    result = Recognise_2({"frame": frame})

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
