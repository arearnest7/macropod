package function

import (
    pb "app/macropod_pb"

    "fmt"
    "context"
    "time"
    "strconv"
    "os"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    structpb "google.golang.org/protobuf/types/known/structpb"
)

type Context struct {
    Function string
    Text string
    JSON map[string]interface{}
    Data []byte
    WorkflowID string
    Depth int32
    Width int32
}

func Invoke(ctx_in Context, dest string, payload string) (string, int32) {
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.WorkflowID + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + ",invoke_rpc_start")
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    depth := ctx_in.Depth + 1
    width := int32(0)
    req := pb.FunctionRequest{Function: &dest, Text: &payload, WorkflowID: &ctx_in.WorkflowID, Depth: &depth, Width: &width}
    var res *pb.MacroPodReply
    var err error
    channel, _ := grpc.Dial(os.Getenv(dest), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    switch comm_type := os.Getenv("COMM_TYPE"); comm_type {
        case "direct":
            stub := pb.NewMacroPodFunctionClient(channel)
            res, err = stub.Invoke(ctx, &req)
        case "gateway":
            stub := pb.NewMacroPodIngressClient(channel)
            res, err = stub.FunctionInvoke(ctx, &req)
        default:
            stub := pb.NewMacroPodFunctionClient(channel)
            res, err = stub.Invoke(ctx, &req)
    }
    channel.Close()
    cancel()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.WorkflowID + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + ",invoke_rpc_end")
    return res.GetReply(), res.GetCode()
}

func Invoke_JSON(ctx_in Context, dest string, payload map[string]interface{}) (string, int32) {
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.WorkflowID + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + ",invoke_rpc_start")
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    depth := ctx_in.Depth + 1
    width := int32(0)
    j, _ := structpb.NewStruct(payload)
    req := pb.FunctionRequest{Function: &dest, JSON: j, WorkflowID: &ctx_in.WorkflowID, Depth: &depth, Width: &width}
    var res *pb.MacroPodReply
    var err error
    channel, _ := grpc.Dial(os.Getenv(dest), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    switch comm_type := os.Getenv("COMM_TYPE"); comm_type {
        case "direct":
            stub := pb.NewMacroPodFunctionClient(channel)
            res, err = stub.Invoke(ctx, &req)
        case "gateway":
            stub := pb.NewMacroPodIngressClient(channel)
            res, err = stub.FunctionInvoke(ctx, &req)
        default:
            stub := pb.NewMacroPodFunctionClient(channel)
            res, err = stub.Invoke(ctx, &req)
    }
    channel.Close()
    cancel()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.WorkflowID + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + ",invoke_rpc_end")
    return res.GetReply(), res.GetCode()
}

func Invoke_Data(ctx_in Context, dest string, payload []byte) (string, int32) {
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.WorkflowID + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + ",invoke_rpc_start")
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    depth := ctx_in.Depth + 1
    width := int32(0)
    req := pb.FunctionRequest{Function: &dest, Data: payload, WorkflowID: &ctx_in.WorkflowID, Depth: &depth, Width: &width}
    var res *pb.MacroPodReply
    var err error
    channel, _ := grpc.Dial(os.Getenv(dest), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    switch comm_type := os.Getenv("COMM_TYPE"); comm_type {
        case "direct":
            stub := pb.NewMacroPodFunctionClient(channel)
            res, err = stub.Invoke(ctx, &req)
        case "gateway":
            stub := pb.NewMacroPodIngressClient(channel)
            res, err = stub.FunctionInvoke(ctx, &req)
        default:
            stub := pb.NewMacroPodFunctionClient(channel)
            res, err = stub.Invoke(ctx, &req)
    }
    channel.Close()
    cancel()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.WorkflowID + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + ",invoke_rpc_end")
    return res.GetReply(), res.GetCode()
}

func Invoke_Multi(ctx_in Context, dest string, payloads []string) ([]string, []int32) {
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.WorkflowID + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + ",invoke_rpc_start")
    tl := make(chan *pb.MacroPodReply)
    for i := 0; i < len(payloads); i++ {
        go func(n int) {
            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            depth := ctx_in.Depth + 1
            width := int32(n)
            req := pb.FunctionRequest{Function: &dest, Text: &payloads[n], WorkflowID: &ctx_in.WorkflowID, Depth: &depth, Width: &width}
            var res *pb.MacroPodReply
            var err error
            channel, _ := grpc.Dial(os.Getenv(dest), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
            switch comm_type := os.Getenv("COMM_TYPE"); comm_type {
                case "direct":
                    stub := pb.NewMacroPodFunctionClient(channel)
                    res, err = stub.Invoke(ctx, &req)
                case "gateway":
                    stub := pb.NewMacroPodIngressClient(channel)
                    res, err = stub.FunctionInvoke(ctx, &req)
                default:
                    stub := pb.NewMacroPodFunctionClient(channel)
                    res, err = stub.Invoke(ctx, &req)
            }
            if err != nil {
                fmt.Println(err)
            }
            tl <- res
            channel.Close()
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
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.WorkflowID + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + ",invoke_rpc_end")
    return reply, code
}

func Invoke_Multi_JSON(ctx_in Context, dest string, payloads []map[string]interface{}) ([]string, []int32) {
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.WorkflowID + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + ",invoke_rpc_start")
    tl := make(chan *pb.MacroPodReply)
    for i := 0; i < len(payloads); i++ {
        go func(n int) {
            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            depth := ctx_in.Depth + 1
            width := int32(n)
            j, _ := structpb.NewStruct(payloads[n])
            req := pb.FunctionRequest{Function: &dest, JSON: j, WorkflowID: &ctx_in.WorkflowID, Depth: &depth, Width: &width}
            var res *pb.MacroPodReply
            var err error
            channel, _ := grpc.Dial(os.Getenv(dest), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
            switch comm_type := os.Getenv("COMM_TYPE"); comm_type {
                case "direct":
                    stub := pb.NewMacroPodFunctionClient(channel)
                    res, err = stub.Invoke(ctx, &req)
                case "gateway":
                    stub := pb.NewMacroPodIngressClient(channel)
                    res, err = stub.FunctionInvoke(ctx, &req)
                default:
                    stub := pb.NewMacroPodFunctionClient(channel)
                    res, err = stub.Invoke(ctx, &req)
            }
            if err != nil {
                fmt.Println(err)
            }
            tl <- res
            channel.Close()
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
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.WorkflowID + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + ",invoke_rpc_end")
    return reply, code
}

func Invoke_Multi_Data(ctx_in Context, dest string, payloads [][]byte) ([]string, []int32) {
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.WorkflowID + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + ",invoke_rpc_start")
    tl := make(chan *pb.MacroPodReply)
    for i := 0; i < len(payloads); i++ {
        go func(n int) {
            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            depth := ctx_in.Depth + 1
            width := int32(n)
            req := pb.FunctionRequest{Function: &dest, Data: payloads[n], WorkflowID: &ctx_in.WorkflowID, Depth: &depth, Width: &width}
            var res *pb.MacroPodReply
            var err error
            channel, _ := grpc.Dial(os.Getenv(dest), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
            switch comm_type := os.Getenv("COMM_TYPE"); comm_type {
                case "direct":
                    stub := pb.NewMacroPodFunctionClient(channel)
                    res, err = stub.Invoke(ctx, &req)
                case "gateway":
                    stub := pb.NewMacroPodIngressClient(channel)
                    res, err = stub.FunctionInvoke(ctx, &req)
                default:
                    stub := pb.NewMacroPodFunctionClient(channel)
                    res, err = stub.Invoke(ctx, &req)
            }
            if err != nil {
                fmt.Println(err)
            }
            tl <- res
            channel.Close()
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
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.WorkflowID + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + ",invoke_rpc_end")
    return reply, code
}
