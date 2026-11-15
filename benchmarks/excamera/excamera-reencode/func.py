import subprocess
import time
import os
import json
import base64

def FunctionHandler(context):
    """handle a request to the function
    Args:
        req (str): request body
    """

    start = time.time() * 1000

    event = context["JSON"]
    req_id = event["req_id"]
    key = event["key"]

    raw_video = f"{key}.y4m"

    prev_state = f"{req_id}-{int(key)-1}-0.state"
    ivf_file = f"{req_id}-{key}-0.ivf"

    interframe_ivf_fn = f"{req_id}-{key}-1.ivf"
    new_state = f"{req_id}-{key}-1.state"

    with open(f"{prev_state}", 'wb') as f:
        f.write(base64.b64decode(event[prev_state].encode("utf-8")))

    with open(f"{ivf_file}", 'wb') as f:
        f.write(base64.b64decode(event[ivf_file].encode("utf-8")))

    # re-encode: replace the first keyframe with an interframe
    #./xc-enc -W -w 0.75 -i y4m -o 2-1.ivf -r -I 1-0.state -p 2-0.ivf -O 2-1.state 2.y4m
    ret = subprocess.run(["./xc-enc",
        "-W",
        "-w",
        "0.75",
        "-i",
        "y4m",
        "-o",
        f'{interframe_ivf_fn}',
        "-r",
        "-I",
        f'{prev_state}',
        "-p"
        f'{ivf_file}',
        "-O",
        f'{new_state}',
        raw_video
        ],
        capture_output=True)

    response = {"req_id": req_id, "time": event["time"]}

    prev_state2 = f"{req_id}-{key}-0.state"

    response[prev_state2] = event[prev_state2]

    with open(f"{interframe_ivf_fn}", "rb") as f: # patching iRoute's broken code (this would break rebase if it did not exist)
        response[interframe_ivf_fn] = base64.b64encode(f.read()).decode("utf-8")

    with open(f"{new_state}", "rb") as f:
        response[new_state] = base64.b64encode(f.read()).decode("utf-8")

    os.remove(f"{prev_state}")
    os.remove(f"{ivf_file}")
    os.remove(f"{interframe_ivf_fn}")
    os.remove(f"{new_state}")

    end = time.time() * 1000

    response["time"][f"reencode-{key}"] = {"start": start, "end": end}

    return json.dumps(response), 200
