from rpc import Invoke_JSON
from rpc import Invoke_Multi_JSON
import json
import uuid
import time
from concurrent.futures import ThreadPoolExecutor

def Entry(req):
    events = []
    for i in range(3):
        event = {}
        event.update(req["req"])
        event["key"] = str(i+1)
        events.append(event)
    all_results_futures, codes = Invoke_Multi_JSON(req["context"], "EXCAMERA_VPXENC", events)
    res = []
    for result in all_results_futures:
        res.append(json.loads(result))
    result_txt, code = Invoke_JSON(req["context"], "EXCAMERA_VX_CON", {"events": res})
    response = json.loads(result_txt)
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
        req = {"req": {"req_id": req_id}, "context": context, "st": st}
        reqs.append(req)

    ex = ThreadPoolExecutor(max_workers=n_req)
    all_results_futures = ex.map(Entry, reqs)
    res = []
    for result in all_results_futures:
        res.append(result)

    return json.dumps(res), 200
