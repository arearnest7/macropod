from parliament import Context
from flask import Request
import json
import os
import random
import time
import traceback
from concurrent import futures

import googlecloudprofiler
from google.auth.exceptions import DefaultCredentialsError
import grpc

import demo_pb2
import demo_pb2_grpc
from grpc_health.v1 import health_pb2
from grpc_health.v1 import health_pb2_grpc

from opentelemetry import trace
from opentelemetry.instrumentation.grpc import GrpcInstrumentorClient, GrpcInstrumentorServer
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter


from logger import getJSONLogger
logger = getJSONLogger('recommendationservice-server')

def ListRecommendations(request):
    max_responses = 5
    # fetch list of products from product catalog stub
    cat_response = os.system("go run productcatalog.go ")
    product_ids = [x.id for x in cat_response.products]
    filtered_products = list(set(product_ids)-set(request.product_ids))
    num_products = len(filtered_products)
    num_return = min(max_responses, num_products)
    # sample list of indicies to return
    indices = random.sample(range(num_products), num_return)
    # fetch product ids from indices
    prod_list = [filtered_products[i] for i in indices]
    logger.info("[Recv ListRecommendations] product_ids={}".format(prod_list))
    # build and return response
    response = demo_pb2.ListRecommendationsResponse()
    response.product_ids.extend(prod_list)
    return response

def main():
    return ListRecommendation(json.loads(argv[1])), 200
