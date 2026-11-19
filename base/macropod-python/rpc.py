from concurrent import futures
import os
import datetime
import random
import grpc
import macropod_pb2 as pb
import macropod_pb2_grpc as pb_grpc
from google.protobuf.struct_pb2 import Struct
from google.protobuf import json_format

MAX_MESSAGE_LENGTH = 1024 * 1024 * 200
opts = [("grpc.max_receive_message_length", MAX_MESSAGE_LENGTH),("grpc.max_send_message_length", MAX_MESSAGE_LENGTH)] # permitting 200MB send and receive maximum
cached_channels = {} # used for cached channel variants and logging

def SendLogger(dest, req):
    channel = None
    if os.environ["LOGGER"] not in cached_channels:
        channel = grpc.insecure_channel(os.environ["LOGGER"], options=opts,)
        cached_channels[os.environ["LOGGER"]] = channel
    else:
        channel = cached_channels[os.environ["LOGGER"]]
    stub = pb_grpc.MacroPodLoggerStub(channel)
    if dest == "timestamp":
        stub.Timestamp(req)
    elif dest == "error":
        stub.Error(req)
    elif dest == "print":
        stub.Print(req)
    else:
        stub.Print(req)

# purpose: send a message to logger (if enabled) to dictate a timestamp of function operation across multiple replicas.
# input:
#   context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
#   target: invocation target for the RPC (if it exists); otherwise null
#   message: message to be stored on the logger if logging is enabled
def Timestamp(context, target, message):
    if "LOGGER" in os.environ and os.environ["LOGGER"] != "":
        req = pb.MacroPodRequest(Workflow=context["Workflow"], Function=context["Function"], WorkflowID=context["WorkflowID"], Depth=context["Depth"], Width=context["Width"], Target=target, Text=message)
        SendLogger("timestamp", req)
    else:
        print(message)

# purpose: send a error to logger (if enabled) notifying that a replica reports an error of err. this is to prevent critical errors for developers/debuggers from being lost due to replica deletion or restart.
# input:
#   context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
#   err: error message to be stored on the logger if logging is enabled
def Error(context, err):
    if "LOGGER" in os.environ and os.environ["LOGGER"] != "":
        req = pb.MacroPodRequest(Workflow=context["Workflow"], Function=context["Function"], Text=err)
        SendLogger("error", req)
    else:
        print(err)

# purpose: send a message to logger (if enabled) to allow persistent printouts for developers/debuggers.
# input:
#   context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
#   message: message to be stored on the logger if logging is enabled
def Print(context, message):
    if "LOGGER" in os.environ and os.environ["LOGGER"] != "":
        req = pb.MacroPodRequest(Workflow=context["Workflow"], Function=context["Function"], Text=message)
        SendLogger("print", req)
    else:
        print(message)

def SendRequests(context, target, reqs):
    res = []
    code = []
    try:
        channel = None
        if "CHANNEL_TYPE" in os.environ and os.environ["CHANNEL_TYPE"] == "cache": # channels are to be cached and reused for future invocations
            if target not in cached_channels:
                channel = grpc.insecure_channel(target, options=opts,)
                cached_channels[target] = channel
            else:
                channel = cached_channels[target]
        elif "CHANNEL_TYPE" in os.environ and os.environ["CHANNEL_TYPE"] == "no-cache": # channels are not to be cached
            channel = grpc.insecure_channel(target, options=opts,)
        else:
            channel = grpc.insecure_channel(target, options=opts,)
    except:
        Error(context, "Invalid destination. Please check your destination for formatting and spelling errors.")
        return [""], [500]
    try:
        if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "direct-w-gateway-fallback-sim" and random.randint(0, 9) == 9:
            fallback_target = target
            ingress_target = os.environ["INGRESS"]
            ingress_channel = None
            if "CHANNEL_TYPE" in os.environ and os.environ["CHANNEL_TYPE"] == "cache": # channels are to be cached and reused for future invocations
                if ingress_target not in cached_channels:
                    ingress_channel = grpc.insecure_channel(ingress_target, options=opts,)
                    cached_channels[ingress_target] = ingress_channel
                else:
                    ingress_channel = cached_channels[ingress_target]
            elif "CHANNEL_TYPE" in os.environ and os.environ["CHANNEL_TYPE"] == "no-cache": # channels are not to be cached
                ingress_channel = grpc.insecure_channel(ingress_target, options=opts,)
            else:
                ingress_channel = grpc.insecure_channel(ingress_target, options=opts,)
            ingress_stub = pb_grpc.MacroPodIngressStub(ingress_channel)
            ingress_req = pb.MacroPodRequest(Workflow=reqs[0].Workflow, Function=reqs[0].Function, Target=reqs[0].Target, WorkflowID=reqs[0].WorkflowID, Depth=reqs[0].Depth, Width=reqs[0].Width)
            fallback_target = ingress_stub.FunctionInvoke(ingress_req)
            if "CHANNEL_TYPE" in os.environ and os.environ["CHANNEL_TYPE"] == "cache": # channels are to be cached and reused for future invocations
                if fallback_target.Reply not in cached_channels:
                    channel = grpc.insecure_channel(fallback_target.Reply, options=opts,)
                    cached_channels[fallback_target.Reply] = channel
                else:
                    channel = cached_channels[fallback_target.Reply]
            elif "CHANNEL_TYPE" in os.environ and os.environ["CHANNEL_TYPE"] == "no-cache": # channels are not to be cached
                channel = grpc.insecure_channel(fallback_target.Reply, options=opts,)
            else:
                channel = grpc.insecure_channel(fallback_target.Reply, options=opts,)
    except:
        Error(context, "Invalid destination. Please check your destination for formatting and spelling errors.")
        return [""], [500]
    tl = []
    with futures.ThreadPoolExecutor(max_workers=len(reqs)) as executor:
        for i in range(len(reqs)):
            if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "direct": # invocation sent directly to target bypassing gateway (direct RPC or IPC)
                stub = pb_grpc.MacroPodFunctionStub(channel)
                tl.append(executor.submit(stub.Invoke, reqs[i]))
            elif "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway": # invocation sent to relay off of gateway (standard serverless behavior)
                stub = pb_grpc.MacroPodIngressStub(channel)
                tl.append(executor.submit(stub.FunctionInvoke, reqs[i]))
            elif "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "direct-w-gateway-fallback-sim":
                stub = pb_grpc.MacroPodFunctionStub(channel)
                tl.append(executor.submit(stub.Invoke, reqs[i]))
            else:
                stub = pb_grpc.MacroPodFunctionStub(channel)
                tl.append(executor.submit(stub.Invoke, reqs[i]))
    try:
        res = [t.result().Reply for t in tl]
        code = [t.result().Code for t in tl]
    except grpc.RpcError as e:
        Error(context, e.details())
    return res, code

# purpose: invoke a downstream function with a text/string payload and return its string output
# input:
#   context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
#   dest: invocation target for the RPC; this is the name of function
#   payload: string payload to be sent to destination function
# output:
#   string: RPC return output
#   int: RPC error code
def Invoke(context, dest, payload):
    Timestamp(context, dest, "Invoke")
    target = os.environ[dest]
    try:
        reqs = [pb.MacroPodRequest(Workflow=context["Workflow"], Function=dest, Target=os.environ[dest], Text=payload, WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=0)]
    except:
        Error(context, "Malformed payload. Please check your formatting. Defaulting to empty message.")
        reqs = [pb.MacroPodRequest(Workflow=context["Workflow"], Function=dest, Target=os.environ[dest], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=0)]
    if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway":
        target = os.environ["INGRESS"]
    res, code = SendRequests(context, target, reqs)
    Timestamp(context, dest, "Invoke_End")
    return res[0], code[0]

# purpose: invoke a downstream function with a JSON list (dictionary) payload and return its string output
# input:
#   context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
#   dest: invocation target for the RPC; this is the name of function
#   payload: dictionary payload to be sent to destination function
# output:
#   string: RPC return output
#   int: RPC error code
def Invoke_JSON(context, dest, payload):
    Timestamp(context, dest, "Invoke_JSON")
    target = os.environ[dest]
    try:
        payload_struct = Struct()
        payload_struct.update(payload)
        reqs = [pb.MacroPodRequest(Workflow=context["Workflow"], Function=dest, Target=os.environ[dest], JSON=payload_struct, WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=0)]
    except:
        Error(context, "Malformed payload. Please check your formatting. Defaulting to empty message.")
        reqs = [pb.MacroPodRequest(Workflow=context["Workflow"], Function=dest, Target=os.environ[dest], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=0)]
    if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway":
        target = os.environ["INGRESS"]
    res, code = SendRequests(context, target, reqs)
    Timestamp(context, dest, "Invoke_JSON_End")
    return res[0], code[0]

# purpose: invoke a downstream function with a bytes payload and return its string output
# input:
#   context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
#   dest: invocation target for the RPC; this is the name of function
#   payload: bytes payload to be sent to destination function
# output:
#   string: RPC return output
#   int: RPC error code
def Invoke_Data(context, dest, payload):
    Timestamp(context, dest, "Invoke_Data")
    target = os.environ[dest]
    try:
        reqs = [pb.MacroPodRequest(Data=payload, Workflow=context["Workflow"], Function=dest, Target=os.environ[dest], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=0)]
    except:
        Error(context, "Malformed payload. Please check your formatting. Defaulting to empty message.")
        reqs = [pb.MacroPodRequest(Workflow=context["Workflow"], Function=dest, Target=os.environ[dest], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=0)]
    if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway":
        target = os.environ["INGRESS"]
    res, code = SendRequests(context, target, reqs)
    Timestamp(context, dest, "Invoke_Data_End")
    return res[0], code[0]

# purpose: invoke a downstream function multiple times concurrently with a list of text/string payloads and return all string outputs
# input:
#   context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
#   dest: invocation target for the RPCs; this is the name of function
#   payloads: list of string payloads to be sent to destination function
# output:
#   string array: array of RPC return outputs
#   int array: array of RPC error code
def Invoke_Multi(context, dest, payloads):
    Timestamp(context, dest, "Invoke_Multi")
    reqs = []
    target = os.environ[dest]
    for i in range(len(payloads)):
        try:
            req = pb.MacroPodRequest(Workflow=context["Workflow"], Function=dest, Target=os.environ[dest], Text=payloads[i], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=i)
        except:
            Error(context, "Malformed payload. Please check your formatting. Defaulting to empty message.")
            req = pb.MacroPodRequest(Workflow=context["Workflow"], Function=dest, Target=os.environ[dest], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=i)
        reqs.append(req)
    if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway":
        target = os.environ["INGRESS"]
    res, code = SendRequests(context, target, reqs)
    Timestamp(context, dest, "Invoke_Multi_End")
    return res, code

# purpose: invoke a downstream function multiple times concurrently with a list of JSON list (dictionary) payloads and return all string outputs
# input:
#   context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
#   dest: invocation target for the RPCs; this is the name of function
#   payloads: list of dictionary payloads to be sent to destination function
# output:
#   string array: array of RPC return outputs
#   int array: array of RPC error code
def Invoke_Multi_JSON(context, dest, payloads):
    Timestamp(context, dest, "Invoke_Multi_JSON")
    reqs = []
    target = os.environ[dest]
    for i in range(len(payloads)):
        try:
            payload_struct = Struct()
            payload_struct.update(payloads[i])
            req = pb.MacroPodRequest(Workflow=context["Workflow"], Function=dest, Target=os.environ[dest], JSON=payload_struct, WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=i)
        except:
            Error(context, "Malformed payload. Please check your formatting. Defaulting to empty message.")
            req = pb.MacroPodRequest(Workflow=context["Workflow"], Function=dest, Target=os.environ[dest], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=i)
        reqs.append(req)
    if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway":
        target = os.environ["INGRESS"]
    res, code = SendRequests(context, target, reqs)
    Timestamp(context, dest, "Invoke_Multi_JSON_End")
    return res, code

# purpose: invoke a downstream function multiple times concurrently with a list of bytes payloads and return all string outputs
# input:
#   context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
#   dest: invocation target for the RPCs; this is the name of function
#   payloads: list of bytes payloads to be sent to destination function
# output:
#   string array: array of RPC return outputs
#   int array: array of RPC error code
def Invoke_Multi_Data(context, dest, payloads):
    Timestamp(context, dest, "Invoke_Multi_Data")
    reqs = []
    target = os.environ[dest]
    for i in range(len(payloads)):
        try:
            req = pb.MacroPodRequest(Data=payloads[i], Workflow=context["Workflow"], Function=dest, Target=os.environ[dest], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=i)
        except:
            Error(context, "Malformed payload. Please check your formatting. Defaulting to empty message.")
            req = pb.MacroPodRequest(Workflow=context["Workflow"], Function=dest, Target=os.environ[dest], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=i)
        reqs.append(req)
    if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway":
        target = os.environ["INGRESS"]
    res, code = SendRequests(context, target, reqs)
    Timestamp(context, dest, "Invoke_Multi_Data_End")
    return res, code
