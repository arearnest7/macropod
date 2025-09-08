from concurrent import futures
import os
import datetime
import grpc
import macropod_pb2 as pb
import macropod_pb2_grpc as pb_grpc
from google.protobuf.struct_pb2 import Struct
from google.protobuf import json_format

MAX_MESSAGE_LENGTH = 1024 * 1024 * 200
opts = [("grpc.max_receive_message_length", MAX_MESSAGE_LENGTH),("grpc.max_send_message_length", MAX_MESSAGE_LENGTH)] # permitting 200MB send and receive maximum

# purpose: send a message to logger (if enabled) to dictate a timestamp of function operation across multiple replicas.
# input:
#   context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
#   target: invocation target for the RPC (if it exists); otherwise null
#   message: message to be stored on the logger if logging is enabled
def Timestamp(context, target, message):
    req = pb.MacroPodRequest(Workflow=context["Workflow"], Function=context["Function"], WorkflowID=context["WorkflowID"], Depth=context["Depth"], Width=context["Width"], Target=target, Text=message)
    if "LOGGER" in os.environ and os.environ["LOGGER"] != "":
        with grpc.insecure_channel(os.environ["LOGGER"], options=opts,) as channel:
            stub = pb_grpc.MacroPodLoggerStub(channel)
            stub.Timestamp(req)

# purpose: send a error to logger (if enabled) notifying that a replica reports an error of err. this is to prevent critical errors for developers/debuggers from being lost due to replica deletion or restart.
# input:
#   context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
#   err: error message to be stored on the logger if logging is enabled
def Error(context, err):
    req = pb.MacroPodRequest(Workflow=context["Workflow"], Function=context["Function"], Text=err)
    if "LOGGER" in os.environ and os.environ["LOGGER"] != "":
        with grpc.insecure_channel(os.environ["LOGGER"], options=opts,) as channel:
            stub = pb_grpc.MacroPodLoggerStub(channel)
            stub.Error(req)

# purpose: send a message to logger (if enabled) to allow persistent printouts for developers/debuggers.
# input:
#   context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
#   message: message to be stored on the logger if logging is enabled
def Print(context, message):
    req = pb.MacroPodRequest(Workflow=context["Workflow"], Function=context["Function"], Text=message)
    if "LOGGER" in os.environ and os.environ["LOGGER"] != "":
        with grpc.insecure_channel(os.environ["LOGGER"], options=opts,) as channel:
            stub = pb_grpc.MacroPodLoggerStub(channel)
            stub.Print(req)

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
    reply = ""
    code = 500
    target = os.environ[dest]
    if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway":
        target = os.environ["INGRESS"]
    req = pb.MacroPodRequest(Workflow=context["Workflow"], Function=dest, Target=os.environ[dest], Text=payload, WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=0)
    with grpc.insecure_channel(target, options=opts,) as channel:
        stub = pb_grpc.MacroPodFunctionStub(channel)
        if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "direct": # invocation sent directly to target bypassing gateway (direct RPC or IPC)
            stub = pb_grpc.MacroPodFunctionStub(channel)
            result = stub.Invoke(req)
            reply = result.Reply
            code = result.Code
        elif "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway": # invocation sent to relay off of gateway (standard serverless behavior)
            stub = pb_grpc.MacroPodIngressStub(channel)
            result = stub.FunctionInvoke(req)
            reply = result.Reply
            code = result.Code
        else:
            result = stub.Invoke(req)
            reply = result.Reply
            code = result.Code
    Timestamp(context, dest, "Invoke_End")
    return reply, code

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
    reply = ""
    code = 500
    target = os.environ[dest]
    if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway":
        target = os.environ["INGRESS"]
    payload_struct = Struct()
    payload_struct.update(payload)
    req = pb.MacroPodRequest(Workflow=context["Workflow"], Function=dest, Target=os.environ[dest], JSON=payload_struct, WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=0)
    with grpc.insecure_channel(target, options=opts,) as channel:
        stub = pb_grpc.MacroPodFunctionStub(channel)
        if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "direct": # invocation sent directly to target bypassing gateway (direct RPC or IPC)
            stub = pb_grpc.MacroPodFunctionStub(channel)
            result = stub.Invoke(req)
            reply = result.Reply
            code = result.Code
        elif "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway": # invocation sent to relay off of gateway (standard serverless behavior)
            stub = pb_grpc.MacroPodIngressStub(channel)
            result = stub.FunctionInvoke(req)
            reply = result.Reply
            code = result.Code
        else:
            result = stub.Invoke(req)
            reply = result.Reply
            code = result.Code
    Timestamp(context, dest, "Invoke_JSON_End")
    return reply, code

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
    reply = ""
    code = 500
    target = os.environ[dest]
    if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway":
        target = os.environ["INGRESS"]
    req = pb.MacroPodRequest(Data=payload, Workflow=context["Workflow"], Function=dest, Target=os.environ[dest], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=0)
    with grpc.insecure_channel(target, options=opts,) as channel:
        stub = pb_grpc.MacroPodFunctionStub(channel)
        if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "direct": # invocation sent directly to target bypassing gateway (direct RPC or IPC)
            stub = pb_grpc.MacroPodFunctionStub(channel)
            result = stub.Invoke(req)
            reply = result.Reply
            code = result.Code
        elif "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway": # invocation sent to relay off of gateway (standard serverless behavior)
            stub = pb_grpc.MacroPodIngressStub(channel)
            result = stub.FunctionInvoke(req)
            reply = result.Reply
            code = result.Code
        else:
            result = stub.Invoke(req)
            reply = result.Reply
            code = result.Code
    Timestamp(context, dest, "Invoke_Data_End")
    return reply, code

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
    reply = []
    code = []
    target = os.environ[dest]
    if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway":
        target = os.environ["INGRESS"]
    with grpc.insecure_channel(target, options=opts,) as channel:
        stub = pb_grpc.MacroPodFunctionStub(channel)
        if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "direct":
            stub = pb_grpc.MacroPodFunctionStub(channel)
        elif "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway":
            stub = pb_grpc.MacroPodIngressStub(channel)
        tl = []
        with futures.ThreadPoolExecutor(max_workers=len(payloads)) as executor:
            for i in range(len(payloads)):
                req = pb.MacroPodRequest(Workflow=context["Workflow"], Function=dest, Target=os.environ[dest], Text=payloads[i], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=i)
                if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "direct": # invocation sent directly to target bypassing gateway (direct RPC or IPC)
                    tl.append(executor.submit(stub.Invoke, req))
                elif "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway": # invocation sent to relay off of gateway (standard serverless behavior)
                    tl.append(executor.submit(stub.FunctionInvoke, req))
                else:
                    tl.append(executor.submit(stub.Invoke, req))
        reply = [t.result().Reply for t in tl]
        code = [t.result().Code for t in tl]
    Timestamp(context, dest, "Invoke_Multi_End")
    return reply, code

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
    reply = []
    code = []
    target = os.environ[dest]
    if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway":
        target = os.environ["INGRESS"]
    with grpc.insecure_channel(target, options=opts,) as channel:
        stub = pb_grpc.MacroPodFunctionStub(channel)
        if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "direct":
            stub = pb_grpc.MacroPodFunctionStub(channel)
        elif "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway":
            stub = pb_grpc.MacroPodIngressStub(channel)
        tl = []
        with futures.ThreadPoolExecutor(max_workers=len(payloads)) as executor:
            for i in range(len(payloads)):
                payload_struct = Struct()
                payload_struct.update(payloads[i])
                req = pb.MacroPodRequest(Workflow=context["Workflow"], Function=dest, Target=os.environ[dest], JSON=payload_struct, WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=i)
                if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "direct": # invocation sent directly to target bypassing gateway (direct RPC or IPC)
                    tl.append(executor.submit(stub.Invoke, req))
                elif "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway": # invocation sent to relay off of gateway (standard serverless behavior)
                    tl.append(executor.submit(stub.FunctionInvoke, req))
                else:
                    tl.append(executor.submit(stub.Invoke, req))
        reply = [t.result().Reply for t in tl]
        code = [t.result().Code for t in tl]
    Timestamp(context, dest, "Invoke_Multi_JSON_End")
    return reply, code

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
    reply = []
    code = []
    target = os.environ[dest]
    if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway":
        target = os.environ["INGRESS"]
    with grpc.insecure_channel(target, options=opts,) as channel:
        stub = pb_grpc.MacroPodFunctionStub(channel)
        if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "direct":
            stub = pb_grpc.MacroPodFunctionStub(channel)
        elif "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway":
            stub = pb_grpc.MacroPodIngressStub(channel)
        tl = []
        with futures.ThreadPoolExecutor(max_workers=len(payloads)) as executor:
            for i in range(len(payloads)):
                req = pb.MacroPodRequest(Data=payloads[i], Workflow=context["Workflow"], Function=dest, Target=os.environ[dest], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=i)
                if "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "direct": # invocation sent directly to target bypassing gateway (direct RPC or IPC)
                    tl.append(executor.submit(stub.Invoke, req))
                elif "COMM_TYPE" in os.environ and os.environ["COMM_TYPE"] == "gateway": # invocation sent to relay off of gateway (standard serverless behavior)
                    tl.append(executor.submit(stub.FunctionInvoke, req))
                else:
                    tl.append(executor.submit(stub.Invoke, req))
        reply = [t.result().Reply for t in tl]
        code = [t.result().Code for t in tl]
    Timestamp(context, dest, "Invoke_Multi_Data_End")
    return reply, code
