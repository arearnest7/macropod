from rpc import RPC
import json
import os
import requests
import base64

def FunctionHandler(context):
    videoFile = open("reference/video.mp4", "rb")
    videoFragment = videoFile.read()
    videoFile.close()
    ret = RPC(context, os.environ['VIDEO_DECODER'], [videoFragment])[0]
    return ret, 200
