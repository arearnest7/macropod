package function

import (
    pb "app/macropod_pb"

    "context"
    "time"
    "os"
    "fmt"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    structpb "google.golang.org/protobuf/types/known/structpb"
)

type Context struct {
    Text       string                   // string payload (if it exists)
    JSON       map[string]interface{}   // dictionary payload (if it exists)
    Data       []byte                   // bytes payload (if it exists)
    Workflow   string                   // workflow name
    Function   string                   // function name
    WorkflowID string                   // invocation id for entire workflow invocation
    Depth      int32                    // invocation depth of current function
    Width      int32                    // invocation entry if multiple invocations occur in parallel (multi RPC); starts at 0 and defaults to 0 if singular invocation
    Target     string                   // internal IP and port of current function
}

var (
    cached_channels = make(map[string]*grpc.ClientConn) // used for cached channel variants and logging
)

func SendLogger(dest string, req pb.MacroPodRequest) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    logger := os.Getenv("LOGGER")
    channel := (*grpc.ClientConn)(nil)
    _, cached := cached_channels[logger]
    if !cached {
        channel, _ := grpc.Dial(logger, grpc.WithDefaultCallOptions(), grpc.WithTransportCredentials(insecure.NewCredentials()))
        cached_channels[logger] = channel
    } else {
        channel = cached_channels[logger]
    }
    channel, _ = grpc.Dial(logger, grpc.WithDefaultCallOptions(), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodLoggerClient(channel)
    switch dest {
        case "timestamp":
            stub.Timestamp(ctx, &req)
        case "error":
            stub.Error(ctx, &req)
        case "print":
            stub.Print(ctx, &req)
        default:
            stub.Print(ctx, &req)
    }
    cancel()
}

// purpose: send a message to logger (if enabled) to dictate a timestamp of function operation across multiple replicas.
// input:
//  context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
//  target: invocation target for the RPC (if it exists); otherwise null
//  message: message to be stored on the logger if logging is enabled
func Timestamp(ctx_in Context, target string, message string) {
    logger := os.Getenv("LOGGER")
    if logger != "" {
        req := pb.MacroPodRequest{Workflow: &ctx_in.Workflow, Function: &ctx_in.Function, WorkflowID: &ctx_in.WorkflowID, Depth: &ctx_in.Depth, Width: &ctx_in.Width, Target: &target, Text: &message}
        SendLogger("timestamp", req)
    } else {
        fmt.Println(message)
    }
}

// purpose: send a error to logger (if enabled) notifying that a replica reports an error of err. this is to prevent critical errors for developers/debuggers from being lost due to replica deletion or restart.
// input:
//  context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
//  err: error message to be stored on the logger if logging is enabled
func Error(ctx_in Context, err error) {
    text := err.Error()
    logger := os.Getenv("LOGGER")
    if logger != "" {
        req := pb.MacroPodRequest{Workflow: &ctx_in.Workflow, Function: &ctx_in.Function, Text: &text}
        SendLogger("error", req)
    } else {
        fmt.Println(text)
    }
}

// purpose: send a message to logger (if enabled) to allow persistent printouts for developers/debuggers.
// input:
//  context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
//  message: message to be stored on the logger if logging is enabled
func Print(ctx_in Context, message string) {
    logger := os.Getenv("LOGGER")
    if logger != "" {
        req := pb.MacroPodRequest{Workflow: &ctx_in.Workflow, Function: &ctx_in.Function, Text: &message}
        SendLogger("print", req)
    } else {
        fmt.Println(message)
    }
}

func SendRequests(ctx_in Context, target string, reqs []pb.MacroPodRequest) ([]string, []int32) {
    var err error
    tl := make(chan *pb.MacroPodReply)
    channel := (*grpc.ClientConn)(nil)
    if os.Getenv("CHANNEL_TYPE") == "cache" {
        _, cached := cached_channels[target]
        if !cached {
            channel, err = grpc.Dial(target, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
            cached_channels[target] = channel
        } else {
            channel = cached_channels[target]
        }
    } else {
        channel, err = grpc.Dial(target, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    }
    if err != nil {
        Error(ctx_in, err)
    }
    for i := 0; i < len(reqs); i++ {
        go func(n int) {
            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            req := reqs[n]
            var res *pb.MacroPodReply
            var err error
            switch comm_type := os.Getenv("COMM_TYPE"); comm_type {
                case "direct": // invocation sent directly to target bypassing gateway (direct RPC or IPC)
                    stub := pb.NewMacroPodFunctionClient(channel)
                    res, err = stub.Invoke(ctx, &req)
                case "gateway": // invocation sent to relay off of gateway (standard serverless behavior)
                    stub := pb.NewMacroPodIngressClient(channel)
                    res, err = stub.FunctionInvoke(ctx, &req)
                default:
                    stub := pb.NewMacroPodFunctionClient(channel)
                    res, err = stub.Invoke(ctx, &req)
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
    for i := 0; i < len(reqs); i++ {
        res := <-tl
        reply = append(reply, res.GetReply())
        code = append(code, res.GetCode())
    }
    if os.Getenv("CHANNEL_TYPE") == "no-cache" {
        channel.Close()
    }
    return reply, code
}

// purpose: invoke a downstream function with a text/string payload and return its string output
// input:
//  context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
//  dest: invocation target for the RPC; this is the name of function
//  payload: string payload to be sent to destination function
// output:
//  string: RPC return output
//  int: RPC error code
func Invoke(ctx_in Context, dest string, payload string) (string, int32) {
    Timestamp(ctx_in, dest, "Invoke")
    depth := ctx_in.Depth + 1
    width := int32(0)
    target := os.Getenv(dest)
    reqs := make([]pb.MacroPodRequest, 0)
    reqs = append(reqs, pb.MacroPodRequest{Workflow: &ctx_in.Workflow, Function: &dest, Target: &target, Text: &payload, WorkflowID: &ctx_in.WorkflowID, Depth: &depth, Width: &width})
    if os.Getenv("COMM_TYPE") == "gateway" {
        target = os.Getenv("INGRESS")
    }
    res, code := SendRequests(ctx_in, target, reqs)
    Timestamp(ctx_in, dest, "Invoke_End")
    return res[0], code[0]
}

// purpose: invoke a downstream function with a JSON list (dictionary) payload and return its string output
// input:
//  context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
//  dest: invocation target for the RPC; this is the name of function
//  payload: dictionary payload to be sent to destination function
// output:
//  string: RPC return output
//  int: RPC error code
func Invoke_JSON(ctx_in Context, dest string, payload map[string]interface{}) (string, int32) {
    Timestamp(ctx_in, dest, "Invoke_JSON")
    depth := ctx_in.Depth + 1
    width := int32(0)
    j, _ := structpb.NewStruct(payload)
    target := os.Getenv(dest)
    reqs := make([]pb.MacroPodRequest, 0)
    reqs = append(reqs, pb.MacroPodRequest{Workflow: &ctx_in.Workflow, Function: &dest, Target: &target, JSON: j, WorkflowID: &ctx_in.WorkflowID, Depth: &depth, Width: &width})
    if os.Getenv("COMM_TYPE") == "gateway" {
        target = os.Getenv("INGRESS")
    }
    res, code := SendRequests(ctx_in, target, reqs)
    Timestamp(ctx_in, dest, "Invoke_JSON_End")
    return res[0], code[0]
}

// purpose: invoke a downstream function with a bytes payload and return its string output
// input:
//  context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
//  dest: invocation target for the RPC; this is the name of function
//  payload: bytes payload to be sent to destination function
// output:
//  string: RPC return output
//  int: RPC error code
func Invoke_Data(ctx_in Context, dest string, payload []byte) (string, int32) {
    Timestamp(ctx_in, dest, "Invoke_Data")
    depth := ctx_in.Depth + 1
    width := int32(0)
    target := os.Getenv(dest)
    reqs := make([]pb.MacroPodRequest, 0)
    reqs = append(reqs, pb.MacroPodRequest{Workflow: &ctx_in.Workflow, Function: &dest, Target: &target, Data: payload, WorkflowID: &ctx_in.WorkflowID, Depth: &depth, Width: &width})
    if os.Getenv("COMM_TYPE") == "gateway" {
        target = os.Getenv("INGRESS")
    }
    res, code := SendRequests(ctx_in, target, reqs)
    Timestamp(ctx_in, dest, "Invoke_Data_End")
    return res[0], code[0]
}

// purpose: invoke a downstream function multiple times concurrently with a list of text/string payloads and return all string outputs
// input:
//  context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
//  dest: invocation target for the RPCs; this is the name of function
//  payloads: list of string payloads to be sent to destination function
// output:
//  string array: array of RPC return outputs
//  int array: array of RPC error code
func Invoke_Multi(ctx_in Context, dest string, payloads []string) ([]string, []int32) {
    Timestamp(ctx_in, dest, "Invoke_Multi")
    reqs := make([]pb.MacroPodRequest, 0)
    target := os.Getenv(dest)
    for i := 0; i < len(payloads); i++ {
        depth := ctx_in.Depth + 1
        width := int32(i)
        req := pb.MacroPodRequest{Workflow: &ctx_in.Workflow, Function: &dest, Target: &target, Text: &payloads[i], WorkflowID: &ctx_in.WorkflowID, Depth: &depth, Width: &width}
        reqs = append(reqs, req)
    }
    if os.Getenv("COMM_TYPE") == "gateway" {
        target = os.Getenv("INGRESS")
    }
    res, code := SendRequests(ctx_in, target, reqs)
    Timestamp(ctx_in, dest, "Invoke_Multi_End")
    return res, code
}

// purpose: invoke a downstream function multiple times concurrently with a list of JSON list (dictionary) payloads and return all string outputs
// input:
//  context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
//  dest: invocation target for the RPCs; this is the name of function
//  payloads: list of dictionary payloads to be sent to destination function
// output:
//  string array: array of RPC return outputs
//  int array: array of RPC error code
func Invoke_Multi_JSON(ctx_in Context, dest string, payloads []map[string]interface{}) ([]string, []int32) {
    Timestamp(ctx_in, dest, "Invoke_Multi_JSON")
    reqs := make([]pb.MacroPodRequest, 0)
    target := os.Getenv(dest)
    for i := 0; i < len(payloads); i++ {
        depth := ctx_in.Depth + 1
        width := int32(i)
        j, _ := structpb.NewStruct(payloads[i])
        req := pb.MacroPodRequest{Workflow: &ctx_in.Workflow, Function: &dest, Target: &target, JSON: j, WorkflowID: &ctx_in.WorkflowID, Depth: &depth, Width: &width}
        reqs = append(reqs, req)
    }
    if os.Getenv("COMM_TYPE") == "gateway" {
        target = os.Getenv("INGRESS")
    }
    res, code := SendRequests(ctx_in, target, reqs)
    Timestamp(ctx_in, dest, "Invoke_Multi_JSON_End")
    return res, code
}

// purpose: invoke a downstream function multiple times concurrently with a list of bytes payloads and return all string outputs
// input:
//  context: context object passed from app handler to function to RPC; this consists of invocation data, including payloads, invocation/workflow id and invocation depth and width
//  dest: invocation target for the RPCs; this is the name of function
//  payloads: list of bytes payloads to be sent to destination function
// output:
//  string array: array of RPC return outputs
//  int array: array of RPC error code
func Invoke_Multi_Data(ctx_in Context, dest string, payloads [][]byte) ([]string, []int32) {
    Timestamp(ctx_in, dest, "Invoke_Multi_Data")
    reqs := make([]pb.MacroPodRequest, 0)
    target := os.Getenv(dest)
    for i := 0; i < len(payloads); i++ {
        depth := ctx_in.Depth + 1
        width := int32(i)
        req := pb.MacroPodRequest{Workflow: &ctx_in.Workflow, Function: &dest, Target: &target, Data: payloads[i], WorkflowID: &ctx_in.WorkflowID, Depth: &depth, Width: &width}
        reqs = append(reqs, req)
    }
    if os.Getenv("COMM_TYPE") == "gateway" {
        target = os.Getenv("INGRESS")
    }
    res, code := SendRequests(ctx_in, target, reqs)
    Timestamp(ctx_in, dest, "Invoke_Multi_Data_End")
    return res, code
}
