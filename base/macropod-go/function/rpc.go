package function

import (
    pb "app/macropod_pb"

    "context"
    "time"
    "os"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    structpb "google.golang.org/protobuf/types/known/structpb"
)

type Context struct {
    Text       string
    JSON       map[string]interface{}
    Data       []byte
    Workflow   string
    Function   string
    WorkflowID string
    Depth      int32
    Width      int32
    Target     string
}

func Timestamp(ctx_in Context, target string, message string) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    req := pb.MacroPodRequest{Workflow: &ctx_in.Workflow, Function: &ctx_in.Function, WorkflowID: &ctx_in.WorkflowID, Depth: &ctx_in.Depth, Width: &ctx_in.Width, Target: &target, Text: &message}
    logger := os.Getenv("LOGGER")
    if logger != "" {
        channel, _ := grpc.Dial(logger, grpc.WithDefaultCallOptions(), grpc.WithTransportCredentials(insecure.NewCredentials()))
        stub := pb.NewMacroPodLoggerClient(channel)
        stub.Timestamp(ctx, &req)
        channel.Close()
        cancel()
    }
}

func Error(ctx_in Context, err error) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    text := err.Error()
    req := pb.MacroPodRequest{Workflow: &ctx_in.Workflow, Function: &ctx_in.Function, Text: &text}
    logger := os.Getenv("LOGGER")
    if logger != "" {
        channel, _ := grpc.Dial(logger, grpc.WithDefaultCallOptions(), grpc.WithTransportCredentials(insecure.NewCredentials()))
        stub := pb.NewMacroPodLoggerClient(channel)
        stub.Error(ctx, &req)
        channel.Close()
        cancel()
    }
}

func Print(ctx_in Context, message string) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    req := pb.MacroPodRequest{Workflow: &ctx_in.Workflow, Function: &ctx_in.Function, Text: &message}
    logger := os.Getenv("LOGGER")
    if logger != "" {
        channel, _ := grpc.Dial(logger, grpc.WithDefaultCallOptions(), grpc.WithTransportCredentials(insecure.NewCredentials()))
        stub := pb.NewMacroPodLoggerClient(channel)
        stub.Print(ctx, &req)
        channel.Close()
        cancel()
    }
}

func Invoke(ctx_in Context, dest string, payload string) (string, int32) {
    Timestamp(ctx_in, dest, "Invoke")
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    depth := ctx_in.Depth + 1
    width := int32(0)
    target := os.Getenv(dest)
    req := pb.MacroPodRequest{Workflow: &ctx_in.Workflow, Function: &dest, Target: &target, Text: &payload, WorkflowID: &ctx_in.WorkflowID, Depth: &depth, Width: &width}
    var res *pb.MacroPodReply
    var err error
    switch comm_type := os.Getenv("COMM_TYPE"); comm_type {
        case "direct":
            channel, _ := grpc.Dial(target, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
            stub := pb.NewMacroPodFunctionClient(channel)
            res, err = stub.Invoke(ctx, &req)
            channel.Close()
        case "gateway":
            channel, _ := grpc.Dial(os.Getenv("INGRESS"), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
            stub := pb.NewMacroPodIngressClient(channel)
            res, err = stub.FunctionInvoke(ctx, &req)
            channel.Close()
        default:
            channel, _ := grpc.Dial(target, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
            stub := pb.NewMacroPodFunctionClient(channel)
            res, err = stub.Invoke(ctx, &req)
            channel.Close()
    }
    cancel()
    if err != nil {
        Error(ctx_in, err)
    }
    Timestamp(ctx_in, dest, "Invoke_End")
    return res.GetReply(), res.GetCode()
}

func Invoke_JSON(ctx_in Context, dest string, payload map[string]interface{}) (string, int32) {
    Timestamp(ctx_in, dest, "Invoke_JSON")
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    depth := ctx_in.Depth + 1
    width := int32(0)
    j, _ := structpb.NewStruct(payload)
    target := os.Getenv(dest)
    req := pb.MacroPodRequest{Workflow: &ctx_in.Workflow, Function: &dest, Target: &target, JSON: j, WorkflowID: &ctx_in.WorkflowID, Depth: &depth, Width: &width}
    var res *pb.MacroPodReply
    var err error
    switch comm_type := os.Getenv("COMM_TYPE"); comm_type {
        case "direct":
            channel, _ := grpc.Dial(target, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
            stub := pb.NewMacroPodFunctionClient(channel)
            res, err = stub.Invoke(ctx, &req)
            channel.Close()
        case "gateway":
            channel, _ := grpc.Dial(os.Getenv("INGRESS"), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
            stub := pb.NewMacroPodIngressClient(channel)
            res, err = stub.FunctionInvoke(ctx, &req)
            channel.Close()
        default:
            channel, _ := grpc.Dial(target, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
            stub := pb.NewMacroPodFunctionClient(channel)
            res, err = stub.Invoke(ctx, &req)
            channel.Close()
    }
    cancel()
    if err != nil {
        Error(ctx_in, err)
    }
    Timestamp(ctx_in, dest, "Invoke_JSON_End")
    return res.GetReply(), res.GetCode()
}

func Invoke_Data(ctx_in Context, dest string, payload []byte) (string, int32) {
    Timestamp(ctx_in, dest, "Invoke_Data")
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    depth := ctx_in.Depth + 1
    width := int32(0)
    target := os.Getenv(dest)
    req := pb.MacroPodRequest{Workflow: &ctx_in.Workflow, Function: &dest, Target: &target, Data: payload, WorkflowID: &ctx_in.WorkflowID, Depth: &depth, Width: &width}
    var res *pb.MacroPodReply
    var err error
    switch comm_type := os.Getenv("COMM_TYPE"); comm_type {
        case "direct":
            channel, _ := grpc.Dial(target, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
            stub := pb.NewMacroPodFunctionClient(channel)
            res, err = stub.Invoke(ctx, &req)
            channel.Close()
        case "gateway":
            channel, _ := grpc.Dial(os.Getenv("INGRESS"), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
            stub := pb.NewMacroPodIngressClient(channel)
            res, err = stub.FunctionInvoke(ctx, &req)
            channel.Close()
        default:
            channel, _ := grpc.Dial(target, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
            stub := pb.NewMacroPodFunctionClient(channel)
            res, err = stub.Invoke(ctx, &req)
            channel.Close()
    }
    cancel()
    if err != nil {
        Error(ctx_in, err)
    }
    Timestamp(ctx_in, dest, "Invoke_Data_End")
    return res.GetReply(), res.GetCode()
}

func Invoke_Multi(ctx_in Context, dest string, payloads []string) ([]string, []int32) {
    Timestamp(ctx_in, dest, "Invoke_Multi")
    tl := make(chan *pb.MacroPodReply)
    for i := 0; i < len(payloads); i++ {
        go func(n int) {
            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            depth := ctx_in.Depth + 1
            width := int32(n)
            target := os.Getenv(dest)
            req := pb.MacroPodRequest{Workflow: &ctx_in.Workflow, Function: &dest, Target: &target, Text: &payloads[n], WorkflowID: &ctx_in.WorkflowID, Depth: &depth, Width: &width}
            var res *pb.MacroPodReply
            var err error
            switch comm_type := os.Getenv("COMM_TYPE"); comm_type {
                case "direct":
                    channel, _ := grpc.Dial(target, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
                    stub := pb.NewMacroPodFunctionClient(channel)
                    res, err = stub.Invoke(ctx, &req)
                    channel.Close()
                case "gateway":
                    channel, _ := grpc.Dial(os.Getenv("INGRESS"), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
                    stub := pb.NewMacroPodIngressClient(channel)
                    res, err = stub.FunctionInvoke(ctx, &req)
                    channel.Close()
                default:
                    channel, _ := grpc.Dial(target, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
                    stub := pb.NewMacroPodFunctionClient(channel)
                    res, err = stub.Invoke(ctx, &req)
                    channel.Close()
            }
            if err != nil {
                Error(ctx_in, err)
            }
            tl <- res
            cancel()
        }(i)
    }
    reply := make([]string, 0)
    code := make([]int32, 0)
    for i := 0; i < len(payloads); i++ {
        res := <-tl
        reply = append(reply, res.GetReply())
        code = append(code, res.GetCode())
    }
    Timestamp(ctx_in, dest, "Invoke_Multi_End")
    return reply, code
}

func Invoke_Multi_JSON(ctx_in Context, dest string, payloads []map[string]interface{}) ([]string, []int32) {
    Timestamp(ctx_in, dest, "Invoke_Multi_JSON")
    tl := make(chan *pb.MacroPodReply)
    for i := 0; i < len(payloads); i++ {
        go func(n int) {
            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            depth := ctx_in.Depth + 1
            width := int32(n)
            j, _ := structpb.NewStruct(payloads[n])
            target := os.Getenv(dest)
            req := pb.MacroPodRequest{Workflow: &ctx_in.Workflow, Function: &dest, Target: &target, JSON: j, WorkflowID: &ctx_in.WorkflowID, Depth: &depth, Width: &width}
            var res *pb.MacroPodReply
            var err error
            switch comm_type := os.Getenv("COMM_TYPE"); comm_type {
                case "direct":
                    channel, _ := grpc.Dial(target, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
                    stub := pb.NewMacroPodFunctionClient(channel)
                    res, err = stub.Invoke(ctx, &req)
                    channel.Close()
                case "gateway":
                    channel, _ := grpc.Dial(os.Getenv("INGRESS"), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
                    stub := pb.NewMacroPodIngressClient(channel)
                    res, err = stub.FunctionInvoke(ctx, &req)
                    channel.Close()
                default:
                    channel, _ := grpc.Dial(target, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
                    stub := pb.NewMacroPodFunctionClient(channel)
                    res, err = stub.Invoke(ctx, &req)
                    channel.Close()
            }
            if err != nil {
                Error(ctx_in, err)
            }
            tl <- res
            cancel()
        }(i)
    }
    reply := make([]string, 0)
    code := make([]int32, 0)
    for i := 0; i < len(payloads); i++ {
        res := <-tl
        reply = append(reply, res.GetReply())
        code = append(code, res.GetCode())
    }
    Timestamp(ctx_in, dest, "Invoke_Multi_JSON_End")
    return reply, code
}

func Invoke_Multi_Data(ctx_in Context, dest string, payloads [][]byte) ([]string, []int32) {
    Timestamp(ctx_in, dest, "Invoke_Multi_Data")
    tl := make(chan *pb.MacroPodReply)
    for i := 0; i < len(payloads); i++ {
        go func(n int) {
            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            depth := ctx_in.Depth + 1
            width := int32(n)
            target := os.Getenv(dest)
            req := pb.MacroPodRequest{Workflow: &ctx_in.Workflow, Function: &dest, Target: &target, Data: payloads[n], WorkflowID: &ctx_in.WorkflowID, Depth: &depth, Width: &width}
            var res *pb.MacroPodReply
            var err error
            switch comm_type := os.Getenv("COMM_TYPE"); comm_type {
                case "direct":
                    channel, _ := grpc.Dial(target, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
                    stub := pb.NewMacroPodFunctionClient(channel)
                    res, err = stub.Invoke(ctx, &req)
                    channel.Close()
                case "gateway":
                    channel, _ := grpc.Dial(os.Getenv("INGRESS"), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
                    stub := pb.NewMacroPodIngressClient(channel)
                    res, err = stub.FunctionInvoke(ctx, &req)
                    channel.Close()
                default:
                    channel, _ := grpc.Dial(target, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
                    stub := pb.NewMacroPodFunctionClient(channel)
                    res, err = stub.Invoke(ctx, &req)
                    channel.Close()
            }
            if err != nil {
                Error(ctx_in, err)
            }
            tl <- res
            cancel()
        }(i)
    }
    reply := make([]string, 0)
    code := make([]int32, 0)
    for i := 0; i < len(payloads); i++ {
        res := <-tl
        reply = append(reply, res.GetReply())
        code = append(code, res.GetCode())
    }
    Timestamp(ctx_in, dest, "Invoke_Multi_Data_End")
    return reply, code
}
