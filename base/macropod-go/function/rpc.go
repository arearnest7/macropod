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
)

type Context struct {
    Function string
    Data []byte
    Text string
    InvocationID string
    Depth int32
    Width int32
}

func Invoke(ctx_in Context, dest string, byte_payloads [][]byte, string_payloads []string) ([]string, []int32) {
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.InvocationID + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + "," + "invoke_rpc_start")
    tl := make(chan *pb.MacroPodReply)
    for i := 0; i < len(byte_payloads); i++ {
        go func(n int) {
            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            depth := ctx_in.Depth + 1
            width := int32(n)
            req := pb.FunctionRequest{Function: dest, Data: byte_payloads[n], InvocationID: &ctx_in.InvocationID, Depth: &depth, Width: &width}
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
    for i := 0; i < len(string_payloads); i++ {
        go func(n int) {
            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            depth := ctx_in.Depth + 1
            width := int32(n)
            req := pb.FunctionRequest{Function: dest, Text: &string_payloads[n], InvocationID: &ctx_in.InvocationID, Depth: &depth, Width: &width}
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
    for i := 0; i < len(byte_payloads) + len(string_payloads); i++ {
        res := <-tl
        reply = append(reply, res.GetReply())
        code = append(code, res.GetCode())
    }
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.InvocationID + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + "," + "invoke_rpc_end")
    return reply, code
}
