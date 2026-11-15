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

    return result_txt

def Entry(req):
    entry_funcs = ["FINRA_MARKETDATA", "FINRA_LASTPX", "FINRA_SIDE", "FINRA_TRDDATE", "FINRA_VOLUME"]
    events = []
    for i in range(5):
        event = req
        events.append({"event": event, "dest": entry_funcs[i]})
    ex = ThreadPoolExecutor(max_workers=5)
    all_results_futures = ex.map(Call, events)
    res = []
    for result in all_results_futures:
        res.append(result)
    headers = {'Content-Type': 'application/json'}
    result_txt = requests.post(os.environ["FINRA_MARGIN_BALANCE"], json=res, headers=headers).text
    response = result_txt
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
        req = {"req_id": req_id, "body": { "portfolioType":"S&P", "portfolio":"1234"}, "st": st}
        reqs.append(req)

    ex = ThreadPoolExecutor(max_workers=n_req)
    all_results_futures = ex.map(Entry, reqs)
    res = []
    for result in all_results_futures:
        res.append(result)

    return json.dumps(res), 200
