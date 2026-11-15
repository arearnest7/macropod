from rpc import Invoke_JSON
import json
import uuid
import time
from concurrent.futures import ThreadPoolExecutor

def Call(event):
    result_txt, code = Invoke_JSON(event["context"], event["dest"], event["event"])

    return result_txt

def Entry(req):
    entry_funcs = ["FINRA_MARKETDATA", "FINRA_LASTPX", "FINRA_SIDE", "FINRA_TRDDATE", "FINRA_VOLUME"]
    events = []
    for i in range(5):
        event = req["req"]
        events.append({"event": event, "dest": entry_funcs[i], "context": req["context"]})
    ex = ThreadPoolExecutor(max_workers=5)
    all_results_futures = ex.map(Call, events)
    res = []
    for result in all_results_futures:
        res.append(result)
    result_txt = Invoke_JSON(req["context"], "FINRA_MARGIN_BALANCE", {"events": res})
    response = result_txt
    time.sleep(req["st"])
    return response

def FunctionHandler(context):
    n_req = 11
    st = 1
    if context["Text"] != "":
        req_j = json.loads(context["Text"])
        n_req = req_j["n"]
        st = req_j["st"]

    reqs = []
    for i in range(n_req):
        req_id = str(uuid.uuid4())
        req = {"req": {"req_id": req_id, "body": {"portfolioType":"S&P", "portfolio":"1234"}}, "context": context, "st": st}
        reqs.append(req)

    ex = ThreadPoolExecutor(max_workers=n_req)
    all_results_futures = ex.map(Entry, reqs)
    res = []
    for result in all_results_futures:
        res.append(result)

    return json.dumps(res), 200
