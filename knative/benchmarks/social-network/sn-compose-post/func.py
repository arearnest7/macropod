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

    return result_txt

def main(context: Context):
    """handle a request to the function
    Args:
        req (str): request body
    """

    start = time.time()

    event = context.request.json

    response = {"time": {"compose-post": {"start_time": start}}}

    try:
        media_ids = (event["media_ids"]).split(",")
        media_types = (event["media_types"]).split(",")
        if not (len(media_ids) == len(media_types)):
            return "Unmatched medias"

        user_id = event["user_id"]
        username = event["username"]
        text = event["text"]
        post_type = event["post_type"]

        response["body"] = event
    except Exception as e:
        response["body"] = "Incomplete arguments"
        return json.dumps(response), 200

    response["time"]["compose-post"]["end_time"] = time.time()

    ex = ThreadPoolExecutor(max_workers=5)
    all_results_futures = ex.map(Call, [{"event": response, "dest": "SN_UPLOAD_USER_MENTIONS"}, {"event": response, "dest": "SN_UPLOAD_CREATOR"}, {"event": response, "dest": "SN_UPLOAD_MEDIA"}, {"event": response, "dest": "SN_UPLOAD_TEXT"}, {"event": response, "dest": "SN_UPLOAD_UNIQUE_ID"}])
    res = []
    for result in all_results_futures:
        res.append(result)
    headers = {'Content-Type': 'application/json'}
    result_txt = requests.post(os.environ["SN_COMPOSE_AND_BUILD"], json=res, headers=headers).text
    response.update(json.loads(result_txt))

    return json.dumps(response), 200
