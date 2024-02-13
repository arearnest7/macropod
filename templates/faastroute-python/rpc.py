from concurrent import futures
import os
import redis
import datetime
import grpc
import app_pb2 as pb
import app_pb2_grpc as pb_grpc

def RPC(dest, payloads, id):
    if "LOGGING_NAME" in os.environ:
        redisClient = redis.Redis(host=os.environ['LOGGING_URL'], password=os.environ['LOGGING_PASSWORD'])
    with grpc.insecure_channel(dest) as channel:
        stub = pb_grpc.gRPCFunctionStub(channel)
        tl = []
        if "LOGGING_NAME" in os.environ:
            redisClient.append(os.environ["LOGGING_NAME"], "INVOKE_START," + id + "," + dest + "," + str(len(payloads)) + "," + str(datetime.datetime.now()) + "\n")
        with futures.ThreadPoolExecutor(max_workers=len(payloads)) as executor:
            for payload in payloads:
                tl.append(executor.submit(stub.gRPCFunctionHandler, pb.RequestBody(body=bytes(payload, "utf-8"), workflow_id=id)))
        if "LOGGING_NAME" in os.environ:
            redisClient.append(os.environ["LOGGING_NAME"], "INVOKE_END," + id + "," + dest + "," + str(len(payloads)) + "," + str(datetime.datetime.now()) + "\n")
        results = [t.result().reply for t in tl]
        return results
