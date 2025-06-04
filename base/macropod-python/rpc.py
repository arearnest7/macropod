from concurrent import futures
import os
import datetime
import grpc
import macropod_pb2 as pb
import macropod_pb2_grpc as pb_grpc

MAX_MESSAGE_LENGTH = 1024 * 1024 * 200
opts = [("grpc.max_receive_message_length", MAX_MESSAGE_LENGTH),("grpc.max_send_message_length", MAX_MESSAGE_LENGTH)]

def Invoke(context, dest, payload):
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + context["WorkflowID"] + "," + str(context["Depth"]) + "," + str(context["Width"]) + ",invoke_rpc_start\n", flush=True)
    reply = ""
    code = 500
    with grpc.insecure_channel(os.environ[dest], options=opts,) as channel:
        stub = pb_grpc.MacroPodFunctionStub(channel)
        if os.environ["COMM_TYPE"] == "direct":
            stub = pb_grpc.MacroPodFunctionStub(channel)
            reply, code = stub.Invoke(pb.RequestBody(Function=dest, Text=payload, WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=0))
        elif os.environ["COMM_TYPE"] == "gateway":
            stub = pb_grpc.MacroPodIngressStub(channel)
            reply, code = stub.FunctionInvoke(pb.RequestBody(Function=dest, Text=payload, WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=0))
        else:
            reply, code = stub.Invoke(pb.RequestBody(Function=dest, Text=payload, WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=0))
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + context["WorkflowID"] + "," + str(context["Depth"]) + "," + str(context["Width"]) + ",invoke_rpc_end\n", flush=True)
    return reply, code

def Invoke_JSON(context, dest, payload):
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + context["WorkflowID"] + "," + str(context["Depth"]) + "," + str(context["Width"]) + ",invoke_rpc_start\n", flush=True)
    reply = ""
    code = 500
    with grpc.insecure_channel(os.environ[dest], options=opts,) as channel:
        stub = pb_grpc.MacroPodFunctionStub(channel)
        if os.environ["COMM_TYPE"] == "direct":
            stub = pb_grpc.MacroPodFunctionStub(channel)
            reply, code = stub.Invoke(pb.RequestBody(Function=dest, JSON=payload, WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=0))
        elif os.environ["COMM_TYPE"] == "gateway":
            stub = pb_grpc.MacroPodIngressStub(channel)
            reply, code = stub.FunctionInvoke(pb.RequestBody(Function=dest, JSON=payload, WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=0))
        else:
            reply, code = stub.Invoke(pb.RequestBody(Function=dest, JSON=payload, WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=0))
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + context["WorkflowID"] + "," + str(context["Depth"]) + "," + str(context["Width"]) + ",invoke_rpc_end\n", flush=True)
    return reply, code

def Invoke_Data(context, dest, payload):
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + context["WorkflowID"] + "," + str(context["Depth"]) + "," + str(context["Width"]) + ",invoke_rpc_start\n", flush=True)
    reply = ""
    code = 500
    with grpc.insecure_channel(os.environ[dest], options=opts,) as channel:
        stub = pb_grpc.MacroPodFunctionStub(channel)
        if os.environ["COMM_TYPE"] == "direct":
            stub = pb_grpc.MacroPodFunctionStub(channel)
            reply, code = stub.Invoke(pb.RequestBody(Function=dest, Data=payload, WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=0))
        elif os.environ["COMM_TYPE"] == "gateway":
            stub = pb_grpc.MacroPodIngressStub(channel)
            reply, code = stub.FunctionInvoke(pb.RequestBody(Function=dest, Data=payload, WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=0))
        else:
            reply, code = stub.Invoke(pb.RequestBody(Function=dest, Data=payload, WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=0))
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + context["WorkflowID"] + "," + str(context["Depth"]) + "," + str(context["Width"]) + ",invoke_rpc_end\n", flush=True)
    return reply, code

def Invoke_Multi(context, dest, payloads):
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + context["WorkflowID"] + "," + str(context["Depth"]) + "," + str(context["Width"]) + ",invoke_rpc_start\n", flush=True)
    reply = []
    code = []
    with grpc.insecure_channel(os.environ[dest], options=opts,) as channel:
        stub = pb_grpc.MacroPodFunctionStub(channel)
        if os.environ["COMM_TYPE"] == "direct":
            stub = pb_grpc.MacroPodFunctionStub(channel)
        elif os.environ["COMM_TYPE"] == "gateway":
            stub = pb_grpc.MacroPodIngressStub(channel)
        tl = []
        with futures.ThreadPoolExecutor(max_workers=len(payloads)) as executor:
            for i in range(len(payloads)):
                if os.environ["COMM_TYPE"] == "direct":
                    tl.append(executor.submit(stub.Invoke, pb.RequestBody(Function=dest, Text=payloads[i], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=i))
                elif os.environ["COMM_TYPE"] == "gateway":
                    tl.append(executor.submit(stub.FunctionInvoke, pb.RequestBody(Function=dest, Text=payloads[i], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=i))
                else:
                    tl.append(executor.submit(stub.Invoke, pb.RequestBody(Function=dest, Text=payloads[i], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=i))
        reply = [t.result().reply for t in tl]
        code = [t.result().code for t in tl]
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + context["WorkflowID"] + "," + str(context["Depth"]) + "," + str(context["Width"]) + ",invoke_rpc_end\n", flush=True)
    return reply, code

def Invoke_Multi_JSON(context, dest, payloads):
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + context["WorkflowID"] + "," + str(context["Depth"]) + "," + str(context["Width"]) + ",invoke_rpc_start\n", flush=True)
    reply = []
    code = []
    with grpc.insecure_channel(os.environ[dest], options=opts,) as channel:
        stub = pb_grpc.MacroPodFunctionStub(channel)
        if os.environ["COMM_TYPE"] == "direct":
            stub = pb_grpc.MacroPodFunctionStub(channel)
        elif os.environ["COMM_TYPE"] == "gateway":
            stub = pb_grpc.MacroPodIngressStub(channel)
        tl = []
        with futures.ThreadPoolExecutor(max_workers=len(payloads)) as executor:
            for i in range(len(payloads)):
                if os.environ["COMM_TYPE"] == "direct":
                    tl.append(executor.submit(stub.Invoke, pb.RequestBody(Function=dest, JSON=payloads[i], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=i))
                elif os.environ["COMM_TYPE"] == "gateway":
                    tl.append(executor.submit(stub.FunctionInvoke, pb.RequestBody(Function=dest, JSON=payloads[i], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=i))
                else:
                    tl.append(executor.submit(stub.Invoke, pb.RequestBody(Function=dest, JSON=payloads[i], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=i))
        reply = [t.result().reply for t in tl]
        code = [t.result().code for t in tl]
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + context["WorkflowID"] + "," + str(context["Depth"]) + "," + str(context["Width"]) + ",invoke_rpc_end\n", flush=True)
    return reply, code

def Invoke_Multi_Data(context, dest, payloads):
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + context["WorkflowID"] + "," + str(context["Depth"]) + "," + str(context["Width"]) + ",invoke_rpc_start\n", flush=True)
    reply = []
    code = []
    with grpc.insecure_channel(os.environ[dest], options=opts,) as channel:
        stub = pb_grpc.MacroPodFunctionStub(channel)
        if os.environ["COMM_TYPE"] == "direct":
            stub = pb_grpc.MacroPodFunctionStub(channel)
        elif os.environ["COMM_TYPE"] == "gateway":
            stub = pb_grpc.MacroPodIngressStub(channel)
        tl = []
        with futures.ThreadPoolExecutor(max_workers=len(payloads)) as executor:
            for i in range(len(payloads)):
                if os.environ["COMM_TYPE"] == "direct":
                    tl.append(executor.submit(stub.Invoke, pb.RequestBody(Function=dest, Data=payloads[i], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=i))
                elif os.environ["COMM_TYPE"] == "gateway":
                    tl.append(executor.submit(stub.FunctionInvoke, pb.RequestBody(Function=dest, Data=payloads[i], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=i))
                else:
                    tl.append(executor.submit(stub.Invoke, pb.RequestBody(Function=dest, Data=payloads[i], WorkflowID=context["WorkflowID"], Depth=(context["Depth"] + 1), Width=i))
        reply = [t.result().reply for t in tl]
        code = [t.result().code for t in tl]
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + context["WorkflowID"] + "," + str(context["Depth"]) + "," + str(context["Width"]) + ",invoke_rpc_end\n", flush=True)
    return reply, code
