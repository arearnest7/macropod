from rpc import RPC
import json
import os
import requests
import base64

def function_handler(context):
    videoFile = open("reference/video.mp4", "rb")
    videoFragment = videoFile.read()
    videoFile.close()
    ret = requests.get(os.environ['VIDEO_DECODER'], [videoFragment.decode()], context["workflow_id"])[0]
    return ret, 200
