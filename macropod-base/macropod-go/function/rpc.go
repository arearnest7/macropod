package function

import (
    "fmt"
    "os"
    "context"
    "time"
    "strconv"
    "math/rand"
    "github.com/go-mmap/mmap"
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
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.WorkflowId + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + "," + ctx_in.RequestType + "," + "3")
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    defer channel.Close()
    stub := pb.NewGRPCFunctionClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    tl := make(chan *pb.ResponseBody)
    pv_paths := make([]string, 0)
    rpc_pv, pv := os.LookupEnv("RPC_PV")
    rpc_dest_pv, dest_pv := os.LookupEnv("RPC_DEST_PV")
    request_type := "gg"
    if !pv {
        if !dest_pv {
            request_type = "gg"
        } else {
            request_type = "gm"
        }
    } else {
        if !dest_pv {
            request_type = "mg"
        } else {
            request_type = "mm"
        }
    }
    for i := 0; i < len(payloads); i++ {
        go func(n int) {
            if request_type == "gg" || request_type == "gm" {
                tl <- invoke(stub, ctx, &pb.RequestBody{Data: payloads[n], WorkflowId: ctx_in.WorkflowId, Depth: int32(ctx_in.Depth + 1), Width: int32(ctx_in.Width), RequestType: &request_type})
            } else {
                pv_path := strconv.Itoa(rand.Intn(10000000))
                pv_paths = append(pv_paths, pv_path)
                os.WriteFile(rpc_pv + "/" + pv_path, payloads[n], 777)
                tl <- invoke(stub, ctx, &pb.RequestBody{WorkflowId: ctx_in.WorkflowId, Depth: int32(ctx_in.Depth + 1), Width: int32(ctx_in.Width), RequestType: &request_type, PvPath: &pv_path})
            }
        }(i)
    }
    results := make([]string, 0)
    for i := 0; i < len(payloads); i++ {
        res := <-tl
        if !dest_pv {
            results = append(results, res.GetReply())
        } else {
            f, _ := mmap.Open(rpc_dest_pv + "/" + res.GetPvPath())
            reply := make([]byte, f.Len())
            f.ReadAt(reply, 0)
            f.Close()
            results = append(results, string(reply))
            os.Remove(rpc_dest_pv + "/" + res.GetPvPath())
        }
    }
    if pv {
        for i := 0; i < len(pv_paths); i++ {
            os.Remove(rpc_pv + "/" + pv_paths[i])
        }
    }
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.WorkflowId + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + "," + ctx_in.RequestType + "," + "4")
    return results
}
