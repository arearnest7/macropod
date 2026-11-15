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

    response["time"]["group"] = {"start": start, "end": end}

    events_re = []
    for i in range(1, 3):
        event = {}
        event.update(response)
        event["key"] = str(i+1)
        events_re.append(event)

    ex = ThreadPoolExecutor(max_workers=2)
    all_results_futures, codes = Invoke_Multi_JSON(context, "EXCAMERA_REENCODE", events_re)
    res = []
    for result in all_results_futures:
        res.append(json.loads(result))
    result_txt, code = Invoke_JSON(context, "EXCAMERA_REBASE", {"events": res})
    response.update(json.loads(result_txt))

    return json.dumps(response), 200
