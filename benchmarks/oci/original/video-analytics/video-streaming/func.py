import json
import os
import requests
import base64

def function_handler(context):
    print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "10" + "\n", flush=True)
    videoFile = open("reference/video.mp4", "rb")
    videoFragment = videoFile.read()
    videoFile.close()
    print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "11" + "\n", flush=True)
    ret = requests.get(os.environ['VIDEO_DECODER'], json={"video": base64.b64encode(videoFragment).decode()}).text
    print(str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "12" + "\n", flush=True)
    return ret, 200
