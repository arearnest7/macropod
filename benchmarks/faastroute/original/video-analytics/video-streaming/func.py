from rpc import RPC
import json
import os
import requests
import base64

def function_handler(context):
    videoFile = open("reference/video.mp4", "rb")
    videoFragment = videoFile.read()
    videoFile.close()
    ret = requests.get(context, os.environ['VIDEO_DECODER'], [videoFragment.decode()])[0]
    return ret, 200
