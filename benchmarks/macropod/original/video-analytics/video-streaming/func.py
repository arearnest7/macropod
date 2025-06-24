from rpc import Invoke_Data
import json
import os
import base64

def FunctionHandler(context):
    videoFile = open("reference/" + context["JSON"]["video"], "rb")
    videoFragment = videoFile.read()
    videoFile.close()
    ret = Invoke_Data(context, 'VIDEO_DECODER', videoFragment)[0]
    return ret, 200
