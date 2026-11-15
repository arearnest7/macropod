import subprocess
import time
import json
import os
import base64

def FunctionHandler(context):
    """handle a request to the function
    Args:
        req (str): request body
    """

    start = time.time() * 1000

    events = context["JSON"]["events"]
    req_id = events[0]["req_id"]

    event = {}
    for event_item in events:
        event.update(event_item)

    raw_video = "3.y4m"
    prev_state = f"{req_id}-2-1.state"
    interframe_ivf_fn = f"{req_id}-3-1.ivf"
    prev_initial_state = f"{req_id}-2-0.state"

    final_state = f"{req_id}-3-1.state"
    final_fn = f"{req_id}-3.ivf"

    with open(f"{prev_state}", 'wb') as f:
        f.write(base64.b64decode(event[prev_state].encode("utf-8")))

    with open(f"{interframe_ivf_fn}", 'wb') as f:
        f.write(base64.b64decode(event[interframe_ivf_fn].encode("utf-8")))

    with open(f"{prev_initial_state}", 'wb') as f:
        f.write(base64.b64decode(event[prev_initial_state].encode("utf-8")))

    # rebase: Without recoding any frames, update my interframe-only ivf file's decoder states
    #./xc-enc -W -w 0.75 -i y4m -o 3.ivf -r -I 2-1.state -p 3-1.ivf -S 2-0.state -O 3-1.state 3.y4m
    ret = subprocess.run(["./xc-enc",
        "-W",
        "-w",
        "0.75",
        "-i",
        "y4m",
        "-o",
        f'{final_fn}',
        "-r",
        "-I",
        f'{prev_state}',
        "-p"
        f'{interframe_ivf_fn}',
        "-S",
        f"{prev_initial_state}",
        "-O",
        f'{final_state}',
        raw_video
        ],
        capture_output=True)

    os.remove(f"{prev_state}")
    os.remove(f"{interframe_ivf_fn}")
    os.remove(f"{prev_initial_state}")
    os.remove(f"{final_state}")
    os.remove(f"{final_fn}")

    response = {"req_id": req_id, "time": {}}

    for event_item in events:
        response["time"].update(event_item["time"])

    end = time.time() * 1000

    response["time"]["rebase"] = {"start": start, "end": end}

    return json.dumps(response), 200
