import json
import os
import requests
import base64
import datetime

def function_handler(context):
    workflow_id = context["request"]["workflow_id"]
    workflow_depth = context["request"]["workflow_depth"]
    workflow_width = context["request"]["workflow_width"]
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "3" + "\n", flush=True)
    videoFile = open("reference/video.mp4", "rb")
    videoFragment = videoFile.read()
    videoFile.close()
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "4" + "\n", flush=True)
    ret = requests.post(os.environ['VIDEO_DECODER'], json={"video": base64.b64encode(videoFragment).decode(), "workflow_id": workflow_id, "workflow_depth": workflow_depth + 1, "workflow_width": 0}).text
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + workflow_id + "," + str(workflow_depth) + "," + str(workflow_width) + "," + "HTTP" + "," + "5" + "\n", flush=True)
    return ret, 200
