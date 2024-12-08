from concurrent import futures
import datetime
import requests
import json

def RPC(context, dest, payloads):
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + context["WorkflowId"] + "," + str(context["Depth"]) + "," + str(context["Width"]) + "," + context["RequestType"] + "," + "rpc_start" + "\n", flush=True)
    request_type = "gg"
    tl = []
    with futures.ThreadPoolExecutor(max_workers=len(payloads)) as executor:
        for i in range(len(payloads)):
            payload = payloads[i]
            tl.append(executor.submit(requests.post, url="http://" + dest, json={"Data": payload.decode(), "WorkflowId": context["WorkflowId"], "Depth": (context["Depth"] + 1), "Width": i, "RequestType": request_type}))
    results = []
    results = [t.result().text for t in tl]
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + context["WorkflowId"] + "," + str(context["Depth"]) + "," + str(context["Width"]) + "," + context["RequestType"] + "," + "rpc_end" + "\n", flush=True)
    return results
