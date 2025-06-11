from __future__ import print_function
from parliament import Context
from flask import Request
import json
from torchvision import transforms
from PIL import Image
import torch
import torchvision.models as models

import pickle
import sys
import os
import io

import cv2
import tempfile
import argparse
import socket
import requests

from concurrent.futures import ThreadPoolExecutor
import base64

import datetime
import random

# Load model
model = models.squeezenet1_1(pretrained=True)

labels_fd = open('imagenet_labels.txt', 'r')
labels = []
for i in labels_fd:
    labels.append(i)
labels_fd.close()

def preprocessImage(imageBytes):
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
    # Set up model to do evaluation
    model.eval()

    # Run inference
    with torch.no_grad():
        out = model(batch_t)

    # Print top 5 for logging
    _, indices = torch.sort(out, descending=True)
    percentages = torch.nn.functional.softmax(out, dim=1)[0] * 100

    # make comma-seperated output of top 100 label
    out = ""
    for idx in indices[0][:100]:
        out = out + labels[idx] + ","
    return out

def Recognise(request):
    # get the frame from s3 or inline
    frame = request

    classification = infer(preprocessImage(frame))
    return classification

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

def Recognise2(frame):
    result = Recognise(frame)

    return result

def processFrames(videoBytes):
    frames = decode(videoBytes)
    all_result_futures = []
    # send all requests
    frames = frames[0:6]
    ex = ThreadPoolExecutor(max_workers=6)
    all_result_futures = ex.map(Recognise2, frames)
    results = ""
    for result in all_result_futures:
        results = results + result + ","

    return results

def Decode(request):
    videoBytes = request
    results = processFrames(videoBytes)
    return results

def main(context: Context):
    if 'request' in context.keys():
        videoFile = open("reference/" + context.request.json["video"], "rb")
        videoFragment = videoFile.read()
        videoFile.close()
        ret = Decode(videoFragment)
        return ret, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
