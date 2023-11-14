from parliament import Context
from flask import Request
import json
import requests
import base64

def main(context: Context):
    if 'request' in context.keys():
        videoFile = open("reference/video.mp4", "rb")
        videoFragment = videoFile.read()
        videoFile.close()
        ret = requests.get(os.environ['VIDEO_DECODER'] + ":80", json={"video": str(bas64.b64encode(videoFragment))}).text
        return ret, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
