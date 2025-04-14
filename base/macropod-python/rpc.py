from concurrent import futures
import os
import datetime
import grpc
import macropod_pb2 as pb
import macropod_pb2_grpc as pb_grpc

MAX_MESSAGE_LENGTH = 1024 * 1024 * 200
opts = [("grpc.max_receive_message_length", MAX_MESSAGE_LENGTH),("grpc.max_send_message_length", MAX_MESSAGE_LENGTH)]

def Invoke(context, dest, byte_payloads, string_payloads):
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + context["InvocationID"] + "," + str(context["Depth"]) + "," + str(context["Width"]) + "," + "invoke_rpc_start" + "\n", flush=True)
    with grpc.insecure_channel(os.environ[dest], options=opts,) as channel:
        stub = pb_grpc.MacroPodFunctionStub(channel)
        if os.environ["COMM_TYPE"] == "direct":
            stub = pb_grpc.MacroPodFunctionStub(channel)
        elif os.environ["COMM_TYPE"] == "gateway":
            stub = pb_grpc.MacroPodIngressStub(channel)
        tl = []
        with futures.ThreadPoolExecutor(max_workers=len(byte_payloads)+len(string_payloads)) as executor:
            for i in range(len(byte_payloads)):
                payload = byte_payloads[i]
                if os.environ["COMM_TYPE"] == "direct":
                    tl.append(executor.submit(stub.Invoke, pb.RequestBody(Function=dest, Data=payload, InvocationID=context["InvocationID"], Depth=(context["Depth"] + 1), Width=i))
                elif os.environ["COMM_TYPE"] == "gateway":
                    tl.append(executor.submit(stub.FunctionInvoke, pb.RequestBody(Function=dest, Data=payload, InvocationID=context["InvocationID"], Depth=(context["Depth"] + 1), Width=i))
                else:
                    tl.append(executor.submit(stub.Invoke, pb.RequestBody(Function=dest, Data=payload, InvocationID=context["InvocationID"], Depth=(context["Depth"] + 1), Width=i))
            for i in range(len(string_payloads)):
                payload = string_payloads[i+len(byte_payloads)]
                if os.environ["COMM_TYPE"] == "direct":
                    tl.append(executor.submit(stub.Invoke, pb.RequestBody(Function=dest, Data=payload, InvocationID=context["InvocationID"], Depth=(context["Depth"] + 1), Width=i+len(byte_payloads)))
                elif os.environ["COMM_TYPE"] == "gateway":
                    tl.append(executor.submit(stub.FunctionInvoke, pb.RequestBody(Function=dest, Data=payload, InvocationID=context["InvocationID"], Depth=(context["Depth"] + 1), Width=i+len(byte_payloads)))
                else:
                    tl.append(executor.submit(stub.Invoke, pb.RequestBody(Function=dest, Data=payload, InvocationID=context["InvocationID"], Depth=(context["Depth"] + 1), Width=i+len(byte_payloads)))

        replies = [t.result().reply for t in tl]
        codes = [t.results().code for t in tl]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + context["InvocationID"] + "," + str(context["Depth"]) + "," + str(context["Width"]) + "," + "invoke_rpc_end" + "\n", flush=True)
        return results, codes
