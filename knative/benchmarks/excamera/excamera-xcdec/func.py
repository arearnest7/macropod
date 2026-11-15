from parliament import Context
from flask import Request
import subprocess
import time
import os
import json
import base64

def main(context: Context):
    """handle a request to the function
    Args:
        req (str): request body
    """

    start = time.time() * 1000

    event = context.request.json
    req_id = event["req_id"]
    key = event["key"]

    vpxenc_fn = f"{req_id}-{key}-0.ivf"
    output_fn = f"{req_id}-{key}-0.state"

    with open(f"{vpxenc_fn}", 'wb') as f:
        f.write(base64.b64decode(event[vpxenc_fn].encode("utf-8")))

    ret = subprocess.run(["./xc-dump",
        f"{vpxenc_fn}",
        f"{output_fn}"],
        capture_output=True)

    response = {"req_id": req_id, "time": event["time"]}

    with open(f"{output_fn}", "rb") as f:
        response[output_fn] = base64.b64encode(f.read()).decode("utf-8")

    os.remove(f"{vpxenc_fn}")
    os.remove(f"{output_fn}")

    end = time.time() * 1000

    response["time"][f"xcdec-{key}"] = {"start": start, "end": end}

    return json.dumps(response), 200
