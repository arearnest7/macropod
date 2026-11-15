from rpc import Invoke_JSON
from rpc import Invoke_Multi_JSON
import time
import json
from concurrent.futures import ThreadPoolExecutor

def FunctionHandler(context):
    """handle a request to the function
    Args:
        req (str): request body
    """

    start = time.time() * 1000

    events = context["JSON"]["events"]
    req_id = events[0]["req_id"]

    response = {"req_id": req_id, "time": {}}

    for event in events:
        response.update(event)

    for event in events:
        response["time"].update(event["time"])

    end = time.time() * 1000

    response["time"]["vx-con"] = {"start": start, "end": end}

    ex = ThreadPoolExecutor(max_workers=3)
    all_results_futures, codes = Invoke_Multi_JSON(context, "EXCAMERA_XCDEC", events)
    res = []
    for result in all_results_futures:
        res.append(json.loads(result))
    for i in range(3):
        res[i].update(events[i])
    result_txt, code = Invoke_JSON(context, "EXCAMERA_GROUP", {"events": res})
    response.update(json.loads(result_txt))

    return json.dumps(response), 200
