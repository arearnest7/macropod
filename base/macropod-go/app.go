package main

import (
    pb "app/macropod_pb"
    "app/function"

    "os"
    "fmt"
    "strconv"
    "time"
    "net"

    "google.golang.org/grpc"
    "golang.org/x/net/context"
)

type FunctionService struct {
    pb.UnimplementedMacroPodFunctionServer
}

var (

)

func Function_Invoke(function_request pb.FunctionRequest) (string, int, error) {
    var reply string
    var code int
    var err error
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + function_request.GetWorkflowID() + "," + strconv.Itoa(int(function_request.GetDepth())) + "," + strconv.Itoa(int(function_request.GetWidth())) + ",request_start")
    defer func() {
        if err := recover(); err != nil {}
    }()
    reply, code = function.FunctionHandler(function.Context{Function: function_request.GetFunction(), Text: function_request.GetText(), JSON: function_request.GetJSON().AsMap(), Data: function_request.GetData(), WorkflowID: function_request.GetWorkflowID(), Depth: function_request.GetDepth(), Width: function_request.GetWidth()})
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + function_request.GetWorkflowID() + "," + strconv.Itoa(int(function_request.GetDepth())) + "," + strconv.Itoa(int(function_request.GetWidth())) + ",request_end")
    return reply, code, err
}

func (s *FunctionService) Invoke(ctx context.Context, in *pb.FunctionRequest) (*pb.MacroPodReply, error) {
    results, code, err := Function_Invoke(*in)
    c := int32(code)
    res := pb.MacroPodReply{Reply: &results, Code: &c}
    return &res, err
}

func main() {
    service_port := os.Getenv("FUNC_PORT")
    if service_port == "" {
        service_port = "8081"
    }
    l, _ := net.Listen("tcp", ":" + service_port)
    s := grpc.NewServer(grpc.MaxSendMsgSize(1024*1024*200), grpc.MaxRecvMsgSize(1024*1024*200))
    pb.RegisterMacroPodFunctionServer(s, &FunctionService{})

    go s.Serve(l)
}
