from rpc import Invoke_JSON
import json
import time
from concurrent.futures import ThreadPoolExecutor

def Call(event):
    result_txt, code = Invoke_JSON(event["context"], event["dest"], event["event"])

    return result_txt

def FunctionHandler(context):
    """handle a request to the function
    Args:
        req (str): request body
    """

    start = time.time()

    event = context["JSON"]

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
    all_results_futures = ex.map(Call, [{"event": response, "dest": "SN_UPLOAD_USER_MENTIONS", "context": context}, {"event": response, "dest": "SN_UPLOAD_CREATOR", "context": context}, {"event": response, "dest": "SN_UPLOAD_MEDIA", "context": context}, {"event": response, "dest": "SN_UPLOAD_TEXT", "context": context}, {"event": response, "dest": "SN_UPLOAD_UNIQUE_ID", "context": context}])
    res = []
    for result in all_results_futures:
        res.append(result)
    result_txt, code = Invoke_JSON(context, "SN_COMPOSE_AND_BUILD", {"events": res})
    response.update(json.loads(result_txt))

    return json.dumps(response), 200
