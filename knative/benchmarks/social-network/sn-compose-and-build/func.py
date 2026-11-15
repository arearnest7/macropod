from parliament import Context
from flask import Request
import json
import time
import os
from concurrent.futures import ThreadPoolExecutor
import requests

def Call(event):
    headers = {'Content-Type': 'application/json'}
    result_txt = requests.post(os.environ[event["dest"]], json=event["event"], headers=headers).text
    result = json.loads(result_txt)

    return result

def main(context: Context):
    """handle a request to the function
    Args:
        req (str): request body
    """

    start = time.time()

    events = context.request.json
    events = [json.loads(event) for event in events]

    body = {}
    response = {"time": {"compose-and-upload": {"start_time": start}}}

    if len(events) < 5:
        body = 'Incomplete arguments'
    else:
        try:
            for event in events:
                body.update(event["body"])
                response["time"].update(event["time"])
        except Exception as e:
            body = 'Incomplete arguments'

    response["body"] = body
    response["time"]["compose-and-upload"]["end_time"] = time.time()

    ex = ThreadPoolExecutor(max_workers=3)
    all_results_futures = ex.map(Call, [{"event": response, "dest": "SN_UPLOAD_HOME_TIMELINE"}, {"event": response, "dest": "SN_UPLOAD_USER_TIMELINE"}, {"event": response, "dest": "SN_POST_STORAGE"}])
    for result in all_results_futures:
        response.update(result)

    return json.dumps(response), 200
