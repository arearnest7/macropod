from rpc import Invoke_JSON
import json
import time
from concurrent.futures import ThreadPoolExecutor

def Call(event):
    result_txt, code = Invoke_JSON(event["context"], event["dest"], event["event"])
    result = json.loads(result_txt)

    return result

def FunctionHandler(context):
    """handle a request to the function
    Args:
        req (str): request body
    """

    start = time.time()

    events = context["JSON"]["events"]
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
    all_results_futures = ex.map(Call, [{"event": response, "dest": "SN_UPLOAD_HOME_TIMELINE", "context": context}, {"event": response, "dest": "SN_UPLOAD_USER_TIMELINE", "context": context}, {"event": response, "dest": "SN_POST_STORAGE", "context": context}])
    for result in all_results_futures:
        response.update(result)

    return json.dumps(response), 200
