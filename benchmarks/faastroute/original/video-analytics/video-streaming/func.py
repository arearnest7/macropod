from rpc import RPC
import json
import os
import requests
import base64

def function_handler(context):
    videoFile = open("reference/video.mp4", "rb")
    videoFragment = videoFile.read()
    videoFile.close()
    ret = requests.get(os.environ['VIDEO_DECODER'], json={"video": base64.b64encode(videoFragment).decode()}).text
    return ret, 200
