# MIT License
#
# Copyright (c) 2021 EASE lab
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

from __future__ import print_function

import sys
import os
import grpc
import argparse
import boto3
import logging as log
import socket

import numpy as np
import pickle

from concurrent import futures

# adding python tracing and storage sources to the system path
sys.path.insert(0, os.getcwd() + '/../proto/')
sys.path.insert(0, os.getcwd() + '/../../../../utils/tracing/python')
sys.path.insert(0, os.getcwd() + '/../../../../utils/storage/python')
import tracing
from storage import Storage
import stacking_pb2_grpc
import stacking_pb2
import destination as XDTdst
import source as XDTsrc
import utils as XDTutil


parser = argparse.ArgumentParser()
parser.add_argument("-dockerCompose", "--dockerCompose", dest="dockerCompose", default=False, help="Env docker compose")
parser.add_argument("-sp", "--sp", dest="sp", default="80", help="serve port")
parser.add_argument("-zipkin", "--zipkin", dest="zipkinURL",
                    default="http://zipkin.istio-system.svc.cluster.local:9411/api/v2/spans",
                    help="Zipkin endpoint url")

args = parser.parse_args()

if tracing.IsTracingEnabled():
    tracing.initTracer("reducer", url=args.zipkinURL)
    tracing.grpcInstrumentClient()
    tracing.grpcInstrumentServer()

INLINE = "INLINE"
S3 = "S3"
XDT = "XDT"
storageBackend = None

# set aws credentials:
AWS_ID = os.getenv('AWS_ACCESS_KEY', "")
AWS_SECRET = os.getenv('AWS_SECRET_KEY', "")
# set aws bucket name:
BUCKET_NAME = os.getenv('BUCKET_NAME','vhive-stacking')

def get_self_ip():
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    try:
        # doesn't even have to be reachable
        s.connect(('10.255.255.255', 1))
        IP = s.getsockname()[0]
    except Exception:
        IP = '127.0.0.1'
    finally:
        s.close()
    return IP


def model_dispatcher(model_name):
    if model_name == 'LinearSVR':
        return LinearSVR
    elif model_name == 'Lasso':
        return Lasso
    elif model_name == 'LinearRegression':
        return LinearRegression
    elif model_name == 'RandomForestRegressor':
        return RandomForestRegressor
    elif model_name == 'KNeighborsRegressor':
        return KNeighborsRegressor
    elif model_name == 'LogisticRegression':
        return LogisticRegression
    else:
        raise ValueError(f"Model {model_name} not found")


class ReducerServicer(stacking_pb2_grpc.ReducerServicer):
    def __init__(self, transferType, XDTconfig=None):

        self.benchName = BUCKET_NAME
        self.transferType = transferType
        if transferType == S3:
            self.s3_client = boto3.resource(
                service_name='s3',
                region_name=os.getenv("AWS_REGION", 'us-west-1'),
                aws_access_key_id=AWS_ID,
                aws_secret_access_key=AWS_SECRET
            )
        elif transferType == XDT:
            if XDTconfig is None:
                log.fatal("Empty XDT config")
            self.XDTconfig = XDTconfig

    def Reduce(self, request, context):
        log.info(f"Reducer is invoked")

        models = []
        predictions = []

        for i, model_pred_tuple in enumerate(request.model_pred_tuples):
            with tracing.Span(f"Reducer gets model {i} from S3"):
                models.append(pickle.loads(storageBackend.get(model_pred_tuple.model_key)))
                predictions.append(pickle.loads(storageBackend.get(model_pred_tuple.pred_key)))

        meta_features = np.transpose(np.array(predictions))

        meta_features_key = 'meta_features'
        models_key = 'models'
        meta_features_key = storageBackend.put(meta_features_key, pickle.dumps(meta_features))

        models_key = storageBackend.put(models_key, pickle.dumps(models))


        return stacking_pb2.ReduceReply(
            models=b'',
            models_key=models_key,
            meta_features=b'',
            meta_features_key=meta_features_key
        )


def serve():
    transferType = os.getenv('TRANSFER_TYPE', S3)
    if transferType == S3:
        global storageBackend
        storageBackend = Storage(BUCKET_NAME)
        log.info("Using inline or s3 transfers")
        max_workers = int(os.getenv("MAX_SERVER_THREADS", 10))
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=max_workers))
        stacking_pb2_grpc.add_ReducerServicer_to_server(
            ReducerServicer(transferType=transferType), server)
        server.add_insecure_port('[::]:' + args.sp)
        server.start()
        server.wait_for_termination()
    elif transferType == XDT:
        log.fatal("XDT not yet supported")
        XDTconfig = XDTutil.loadConfig()
    else:
        log.fatal("Invalid Transfer type")


if __name__ == '__main__':
    log.basicConfig(level=log.INFO)
    serve()
