# MIT License
#
# Copyright (c) 2021 Michal Baczun and EASE lab
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

"""Azure Function to perform inference.
"""

from torchvision import transforms
from PIL import Image
import torch
import torchvision.models as models

import sys
import os
import pickle
# adding python tracing and storage sources to the system path
sys.path.insert(0, os.getcwd() + '/../proto/')
sys.path.insert(0, os.getcwd() + '/../../../../utils/tracing/python')
sys.path.insert(0, os.getcwd() + '/../../../../utils/storage/python')
import tracing
import videoservice_pb2_grpc
import videoservice_pb2
import destination as XDTdst
import utils as XDTutil

import io
import grpc
import argparse
import boto3
import logging as log

from concurrent import futures


parser = argparse.ArgumentParser()
parser.add_argument("-sp", "--sp", dest = "sp", default = "80", help="serve port")
parser.add_argument("-zipkin", "--zipkin", dest = "url", default = "http://zipkin.istio-system.svc.cluster.local:9411/api/v2/spans", help="Zipkin endpoint url")

args = parser.parse_args()

INLINE = "INLINE"
S3 = "S3"
XDT = "XDT"
storageBackend = None

if tracing.IsTracingEnabled():
    tracing.initTracer("recog", url=args.url)
    tracing.grpcInstrumentServer()

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

class ObjectRecognitionServicer(videoservice_pb2_grpc.ObjectRecognitionServicer):
    def __init__(self, transferType):

        self.transferType = transferType

    def Recognise(self, request, context):
        log.info("received a call")

        # get the frame from s3 or inline
        frame = None
        if self.transferType == S3:
            log.info("retrieving target frame '%s' from s3" % request.s3key)
            with tracing.Span("Frame fetch"):
                frame = pickle.loads(storageBackend.get(request.s3key))
        elif self.transferType == INLINE:
            frame = request.frame

        log.info("performing image recognition on frame")
        classification = infer(preprocessImage(frame))
        log.info("object recogintion successful")
        return videoservice_pb2.RecogniseReply(classification=classification)

def serve():
    transferType = os.getenv('TRANSFER_TYPE', INLINE)
    if transferType == S3:
        from storage import Storage
        bucketName = os.getenv('BUCKET_NAME', 'vhive-video-bench')
        global storageBackend
        storageBackend = Storage(bucketName)
    if transferType == S3 or transferType == INLINE:
        max_workers = int(os.getenv("MAX_RECOG_SERVER_THREADS", 10))
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=max_workers))
        videoservice_pb2_grpc.add_ObjectRecognitionServicer_to_server(
            ObjectRecognitionServicer(transferType=transferType), server)
        server.add_insecure_port('[::]:'+args.sp)
        server.start()
        server.wait_for_termination()
    elif transferType == XDT:
        config = XDTutil.loadConfig()
        log.info("[recog] transfering via XDT")
        log.info(config)

        def handler(imageBytes):
            classification = infer(preprocessImage(imageBytes))
            return classification.encode(), True

        XDTdst.StartDstServer(config, handler)


if __name__ == '__main__':
    log.basicConfig(level=log.INFO)
    serve()
