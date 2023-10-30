from parliament import Context
from flask import Request
import json
from concurrent import futures
import argparse
import os
import sys
import time
import grpc
import traceback
from jinja2 import Environment, FileSystemLoader, select_autoescape, TemplateError
from google.api_core.exceptions import GoogleAPICallError
from google.auth.exceptions import DefaultCredentialsError

import demo_pb2
import demo_pb2_grpc
from grpc_health.v1 import health_pb2
from grpc_health.v1 import health_pb2_grpc

from opentelemetry import trace
from opentelemetry.instrumentation.grpc import GrpcInstrumentorServer
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter

import googlecloudprofiler

from logger import getJSONLogger
logger = getJSONLogger('emailservice-server')

# Loads confirmation email template from file
env = Environment(
    loader=FileSystemLoader('templates'),
    autoescape=select_autoescape(['html', 'xml'])
)
template = env.get_template('confirmation.html')

def send_message(sender, envelope_from_authority, header_from_authority, envelope_from_address, simple_message):
	return {"rfc822_message_id": 1234};
def send_email(email_address, content):
    response = send_message(
        sender = [project_id, region, sender_id],
        envelope_from_authority = '',
        header_from_authority = '',
        envelope_from_address = from_address,
        simple_message = {
            "from": {
            "address_spec": from_address,
            },
            "to": [{
                "address_spec": email_address
            }],
            "subject": "Your Confirmation Email",
            "html_body": content
        }
    )
    logger.info("Message sent: {}".format(response.rfc822_message_id))

def SendOrderConfirmation(request, context):
    email = request.email
    order = request.order

    try:
        confirmation = template.render(order = order)
    except TemplateError as err:
        context.set_details("An error occurred when preparing the confirmation mail.")
        logger.error(err.message)
        context.set_code(grpc.StatusCode.INTERNAL)
        return demo_pb2.Empty()

    try:
        send_email(email, confirmation)
    except GoogleAPICallError as err:
        context.set_details("An error occurred when sending the email.")
        print(err.message)
        context.set_code(grpc.StatusCode.INTERNAL)
        return demo_pb2.Empty()

    return demo_pb2.Empty()

def main(context: Context):
    return SendOrderConfirmation(context.request.json), 200
