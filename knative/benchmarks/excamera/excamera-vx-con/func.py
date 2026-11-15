from parliament import Context
from flask import Request
import time
import json
import os
from concurrent.futures import ThreadPoolExecutor
import requests

def Call(event):
    headers = {'Content-Type': 'application/json'}
    result_txt = requests.post(os.environ["EXCAMERA_XCDEC"], json=event, headers=headers).text
    result = json.loads(result_txt)

    return result

def main(context: Context):
    """handle a request to the function
    Args:
        req (str): request body
    """

    start = time.time() * 1000

    events = context.request.json
    req_id = events[0]["req_id"]

    response = {"req_id": req_id, "time": {}}

    for event in events:
        response.update(event)

    for event in events:
        response["time"].update(event["time"])

    end = time.time() * 1000

    response["time"]["vx-con"] = {"start": start, "end": end}

    ex = ThreadPoolExecutor(max_workers=3)
    all_results_futures = ex.map(Call, events)
    res = []
    for result in all_results_futures:
        res.append(result)
    for i in range(3):
        res[i].update(events[i])
    headers = {'Content-Type': 'application/json'}
    result_txt = requests.post(os.environ["EXCAMERA_GROUP"], json=res, headers=headers).text
    response.update(json.loads(result_txt))

    return json.dumps(response), 200
