from concurrent import futures
import datetime
import grpc
import wf_pb2 as pb
import wf_pb2_grpc as pb_grpc

MAX_MESSAGE_LENGTH = 1024 * 1024 * 200
opts = [("grpc.max_receive_message_length", MAX_MESSAGE_LENGTH),("grpc.max_send_message_length", MAX_MESSAGE_LENGTH)]

def RPC(context, dest, payloads):
    print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + context["WorkflowId"] + "," + str(context["Depth"]) + "," + str(context["Width"]) + "," + context["RequestType"] + "," + "rpc_start" + "\n", flush=True)
    with grpc.insecure_channel(dest, options=opts,) as channel:
        stub = pb_grpc.gRPCFunctionStub(channel)
        request_type = "gg"
        with futures.ThreadPoolExecutor(max_workers=len(payloads)) as executor:
            for i in range(len(payloads)):
                payload = payloads[i]
                tl.append(executor.submit(stub.gRPCFunctionHandler, pb.RequestBody(data=payload, workflow_id=context["WorkflowId"], depth=(context["Depth"] + 1), width=i, request_type=request_type)))
        results = []
        results = [t.result().reply for t in tl]
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + context["WorkflowId"] + "," + str(context["Depth"]) + "," + str(context["Width"]) + "," + context["RequestType"] + "," + "rpc_end" + "\n", flush=True)
        return results
