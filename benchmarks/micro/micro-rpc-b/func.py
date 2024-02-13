from parliament import Context
from flask import Request
import json
import os

def main(context: Context):
    return "a" * int(os.environ["LEN"]), 200
