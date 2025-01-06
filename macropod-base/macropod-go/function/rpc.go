package function

import (
    "fmt"
    "context"
    "time"
    "strconv"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    pb "app/wf_pb"
)

type Context struct {
    Request []byte
    WorkflowId string
    Depth int
    Width int
    RequestType string
    InvokeType string
    IsJson bool
}

func invoke(stub pb.GRPCFunctionClient, ctx_in context.Context, in *pb.RequestBody) (*pb.ResponseBody) {
    res, _ := stub.GRPCFunctionHandler(ctx_in, in)
    return res
}

func RPC(ctx_in Context, dest string, payloads [][]byte) ([]string) {
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.WorkflowId + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + "," + ctx_in.RequestType + "," + "rpc_start")
    tl := make(chan *pb.ResponseBody)
    request_type := "gg"
    for i := 0; i < len(payloads); i++ {
        go func(n int) {
            channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
            stub := pb.NewGRPCFunctionClient(channel)
            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            tl <- invoke(stub, ctx, &pb.RequestBody{Data: payloads[n], WorkflowId: ctx_in.WorkflowId, Depth: int32(ctx_in.Depth + 1), Width: int32(i), RequestType: &request_type})
            channel.Close()
            cancel()
        }(i)
    }
    results := make([]string, 0)
    for i := 0; i < len(payloads); i++ {
        res := <-tl
        results = append(results, res.GetReply())
    }
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.WorkflowId + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + "," + ctx_in.RequestType + "," + "rpc_end")
    return results
}
