from concurrent import futures
import os
import datetime
import grpc
import macropod_pb2 as pb
import macropod_pb2_grpc as pb_grpc
from func import FunctionHandler

MAX_MESSAGE_LENGTH = 1024 * 1024 * 200
opts = [("grpc.max_receive_message_length", MAX_MESSAGE_LENGTH),("grpc.max_send_message_length", MAX_MESSAGE_LENGTH)]

class MacroPodFunctionServicer(pb_grpc.MacroPodFunctionServicer):
    def Invoke(self, request, context):
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + request.WorkflowID + "," + str(request.Depth) + "," + str(request.Width) + ",request_start\n", flush=True)
        reply, code = FunctionHandler({"Function": request.Function, "Text": request.Text, "JSON": request.JSON, "Data": request.Data, "WorkflowID": request.WorkflowID, "Depth": request.Depth, "Width": request.Width})
        res = pb.MacroPodReply(Reply=reply, Code=code)
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + request.WorkflowID + "," + str(request.Depth) + "," + str(request.Width) + ",request_end\n", flush=True)
        return res

if __name__ == '__main__':
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=100), options=opts)
    pb_grpc.add_MacroPodFunctionServicer_to_server(MacroPodFunctionServicer(), server)
    server.add_insecure_port("[::]:" + os.environ["FUNC_PORT"])
    server.start()
    server.wait_for_termination()
