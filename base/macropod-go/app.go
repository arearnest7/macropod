package main

import (
    "os"
    "context"
    "fmt"
    "net"
    "time"
    "strconv"
    "google.golang.org/grpc"

    pb "app/wf_pb"
    function "app/function"
)

type server struct {
    pb.UnimplementedGRPCFunctionServer
}

func (s *server) GRPCFunctionHandler(ctx context.Context, in *pb.RequestBody) (*pb.ResponseBody, error) {
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + in.GetWorkflowId() + "," + strconv.Itoa(int(in.GetDepth())) + "," + strconv.Itoa(int(in.GetWidth())) + "," + in.GetRequestType() + "," + "request_start")
    res := pb.ResponseBody{Code: int32(500)}
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + in.GetWorkflowId() + "," + strconv.Itoa(int(in.GetDepth())) + "," + strconv.Itoa(int(in.GetWidth())) + "," + in.GetRequestType() + "," + "function_start")
    reply, code := function.FunctionHandler(function.Context{Request: in.GetData(), WorkflowId: in.GetWorkflowId(), Depth: int(in.GetDepth()), Width: int(in.GetWidth()), RequestType: in.GetRequestType(), InvokeType: "GRPC", IsJson: false})
    res.Reply = &reply
    res.Code = int32(code)
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + in.GetWorkflowId() + "," + strconv.Itoa(int(in.GetDepth())) + "," + strconv.Itoa(int(in.GetWidth())) + "," + in.GetRequestType() + "," + "request_end")
    return &res, nil
}

func main() {
    port, _ := strconv.Atoi(os.Getenv("FUNC_PORT"))
    l, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))
    s := grpc.NewServer(grpc.MaxSendMsgSize(1024*1024*200), grpc.MaxRecvMsgSize(1024*1024*200))
    pb.RegisterGRPCFunctionServer(s, &server{})
    s.Serve(l)
}
