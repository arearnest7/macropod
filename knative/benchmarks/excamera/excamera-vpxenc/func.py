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

    raw_video = f"{key}.y4m"

    vpxenc_fn = f"{req_id}-{key}-vpxenc.ivf"
    output_fn = f"{req_id}-{key}-0.ivf"

    ret = subprocess.run(["./vpxenc",
        "--ivf",
        "--codec=vp8",
        "--good",
        "--cpu-used=0",
        "--end-usage=cq",
        "--min-q=0",
        "--max-q=63",
        "--cq-level=22",
        "--buf-initial-sz=10000",
        "--buf-optimal-sz=20000",
        "--buf-sz=40000",
        "--undershoot-pct=100",
        "--passes=2",
        "--auto-alt-ref=1",
        "--threads=1",
        "--token-parts=0",
        "--tune=ssim",
        "--target-bitrate=4294967295",
        "-o",
        f"{vpxenc_fn}",
        raw_video],
        capture_output=True)

    ret = subprocess.run(["./xc-terminate-chunk",
        f"{vpxenc_fn}",
        f"{output_fn}"],
        capture_output=True)

    with open(f"{output_fn}", "rb") as f:
        event[output_fn] = base64.b64encode(f.read()).decode("utf-8")

    os.remove(f"{vpxenc_fn}")
    os.remove(f"{output_fn}")

    end = time.time() * 1000

    event["time"] = {f"vpxenc-{key}": {"start": start, "end": end}}

    return json.dumps(event), 200
