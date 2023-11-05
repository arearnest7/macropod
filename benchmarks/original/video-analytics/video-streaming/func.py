from parliament import Context
from flask import Request
import json
import logging
import requests

def main(context: Context):
    if 'request' in context.keys():
        videoFile = open("reference/video.mp4", "rb")
        logging.info('read video fragment, size: {}\n'.format(len(videoFragment)))
        videoFragment = videoFile.read()
        videoFile.close()
        ret = requests.get(url = os.environ['VIDEO_DECODER'] + ":8080", data = json.dumps({"video": videoFragment})).text
        return ret, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
