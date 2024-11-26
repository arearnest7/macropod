from concurrent import futures
import os
import datetime
import grpc
import wf_pb2 as pb
import wf_pb2_grpc as pb_grpc
from func import FunctionHandler

app = Flask(__name__)
api = Api(app)

MAX_MESSAGE_LENGTH = 1024 * 1024 * 200
opts = [("grpc.max_receive_message_length", MAX_MESSAGE_LENGTH),("grpc.max_send_message_length", MAX_MESSAGE_LENGTH)]

class gRPCFunctionServicer(pb_grpc.gRPCFunctionServicer):
    def gRPCFunctionHandler(self, request, context):
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + request.workflow_id + "," + str(request.depth) + "," + str(request.width) + "," + request.request_type + "," + "request_start" + "\n", flush=True)
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + request.workflow_id + "," + str(request.depth) + "," + str(request.width) + "," + request.request_type + "," + "function_start" + "\n", flush=True)
        reply, code = FunctionHandler({"Request": request.data, "WorkflowId": request.workflow_id, "Depth": request.depth, "Width": request.width, "RequestType": request.request_type, "InvokeType": "GRPC", "IsJson": False})
        res = pb.ResponseBody(reply=reply, code=code)
        print(datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f %Z") + "," + request.workflow_id + "," + str(request.depth) + "," + str(request.width) + "," + request.request_type + "," + "request_end" + "\n", flush=True)
        return res

if __name__ == '__main__':
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=int(os.environ["GRPC_THREAD"])), options=opts)
    pb_grpc.add_gRPCFunctionServicer_to_server(gRPCFunctionServicer(), server)
    server.add_insecure_port("[::]:" + os.environ["FUNC_PORT"])
    server.start()
    server.wait_for_termination()
