from parliament import Context
from flask import Request
import base64
import random
import os
import datetime
import redis


def main(context: Context):
    return "CategorizationFail: Fail: \"Input CSV could not be categorised into 'Product' or 'Service'.\"", 200
