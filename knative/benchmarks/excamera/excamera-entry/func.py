from parliament import Context
from flask import Request
import json
import uuid
import time
import os
from concurrent.futures import ThreadPoolExecutor
import requests

def Call(event):
    headers = {'Content-Type': 'application/json'}
    result_txt = requests.post(os.environ[event["dest"]], json=event["event"], headers=headers).text
    result = json.loads(result_txt)

    return result

def Entry(req):
    events = []
    for i in range(3):
        event = {}
        event.update(req["req"])
        event["key"] = i+1
        events.append({"event": event, "dest": "EXCAMERA_VPXENC"})
    ex = ThreadPoolExecutor(max_workers=3)
    all_results_futures = ex.map(Call, events)
    res = []
    for result in all_results_futures:
        res.append(result)
    headers = {'Content-Type': 'application/json'}
    result_txt = requests.post(os.environ["EXCAMERA_VX_CON"], json=res, headers=headers).text
    response = json.loads(result_txt)
    time.sleep(req["st"])
    return response

def main(context: Context):
    n_req = 11
    st = 1
    if 'request' in context.keys():
        req_j = json.loads(context.request.get_data(as_text=True))
        n_req = req_j["n"]
        st = req_j["st"]

    reqs = []
    for i in range(n_req):
        req_id = str(uuid.uuid4())
        req = {"req": {"req_id": req_id}, "st": st}
        reqs.append(req)

    ex = ThreadPoolExecutor(max_workers=n_req)
    all_results_futures = ex.map(Entry, reqs)
    res = []
    for result in all_results_futures:
        res.append(result)

    return json.dumps(res), 200
