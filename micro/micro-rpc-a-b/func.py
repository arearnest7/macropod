from parliament import Context
from flask import Request
import json
import os
import requests

def b_main(a):
    return "a" * int(os.environ["LEN_B"])

def main(context: Context):
    return b_main("a" * int(os.environ["LEN_A"])), 200
