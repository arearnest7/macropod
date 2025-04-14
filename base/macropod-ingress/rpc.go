package app

import (
    pb "app/macropod_pb"

    "fmt"
    "context"
    "time"
    "strconv"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

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

func DeployerConfig(config pb.ConfigStruct, dest string) {
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodDeployerClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    _, err := stub.Config(ctx, &config)
    if err != nil {
        fmt.Println(err)
    }
}

func DeployerCreateWorkflow(workflow pb.WorkflowStruct, dest string) {
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodDeployerClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    _, err := stub.CreateWorkflow(ctx, &workflow)
    if err != nil {
        fmt.Println(err)
    }
}

func DeployerUpdateWorkflow(workflow pb.WorkflowStruct, dest string) {
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodDeployerClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    _, err := stub.UpdateWorkflow(ctx, &workflow)
    if err != nil {
        fmt.Println(err)
    }
}

func DeployerDeleteWorkflow(deployer_request pb.DeployerRequest, dest string) {
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodDeployerClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    _, err := stub.DeleteWorkflow(ctx, &deployer_request)
    if err != nil {
        fmt.Println(err)
    }
}

func DeployerUpdateDeployments(deployer_request pb.DeployerRequest, dest string) (int) {
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodDeployerClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    res, err := stub.UpdateDeployments(ctx, &deployer_request)
    if err != nil {
        fmt.Println(err)
    }
    return int(res.GetCode())
}

func DeployerCreateDeployment(deployer_request pb.DeployerRequest, dest string) (int) {
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodDeployerClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    res, err := stub.CreateDeployment(ctx, &deployer_request)
    if err != nil {
        fmt.Println(err)
    }
    return int(res.GetCode())
}

func DeployerTTLDelete(deployer_request pb.DeployerRequest, dest string) {
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodDeployerClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    _, err := stub.TTLDelete(ctx, &deployer_request)
    if err != nil {
        fmt.Println(err)
    }
}

func DeployerUpdateManifest(manifest pb.MacroPodManifest, dest string) (pb.MacroPodManifest) {
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodDeployerClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    deployer_manifest, err := stub.UpdateManifest(ctx, &manifest)
    if err != nil {
        fmt.Println(err)
    }
    return &deployer_manifest
}

func IngressConfig(config pb.ConfigStruct, dest string) {
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodIngressClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    _, err := stub.Config(ctx, &config)
    if err != nil {
        fmt.Println(err)
    }
}

func IngressWorkflowInvoke(ingress_request pb.IngressRequest, dest string) (string, int) {
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodIngressClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    res, err := stub.WorkflowInvoke(ctx, &ingress_request)
    if err != nil {
        fmt.Println(err)
    }
    return res.GetReply(), int(res.GetCode())
}

func IngressFunctionInvoke(function_request pb.FunctionRequest, dest string) (string, int) {
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodIngressClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    res, err := stub.FunctionInvoke(ctx, &function_request)
    if err != nil {
        fmt.Println(err)
    }
    return res.GetReply(), int(res.GetCode())
}

func IngressCreateWorkflow(workflow pb.WorkflowStruct, dest string) {
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodIngressClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    _, err := stub.CreateWorkflow(ctx, &workflow)
    if err != nil {
        fmt.Println(err)
    }
}

func IngressUpdateWorkflow(workflow pb.WorkflowStruct, dest string) {
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodIngressClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    _, err := stub.UpdateWorkflow(ctx, &workflow)
    if err != nil {
        fmt.Println(err)
    }
}

func IngressDeleteWorkflow(ingress_request pb.IngressRequest, dest string) {
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodIngressClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    _, err := stub.DeleteWorkflow(ctx, &ingress_request)
    if err != nil {
        fmt.Println(err)
    }
}

func IngressEval(workflow pb.WorkflowStruct, dest string) (string) {
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodIngressClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    res, err := stub.Eval(ctx, &workflow)
    if err != nil {
        fmt.Println(err)
    }
    return res.GetReply()
}

func IngressEvalMetrics(ingress_eval_request pb.IngressEvalRequest, dest string) (string) {
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodIngressClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    res, err := stub.EvalMetrics(ctx, &ingress_eval_request)
    if err != nil {
        fmt.Println(err)
    }
    return res.GetReply()
}

func IngressEvalLatency(ingress_eval_request pb.IngressEvalRequest) (string) {
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodIngressClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    res, err := stub.EvalLatency(ctx, &ingress_eval_request)
    if err != nil {
        fmt.Println(err)
    }
    return res.GetReply()
}

func IngressEvalSummary(ingress_eval_request pb.IngressEvalRequest) (string) {
    channel, _ := grpc.Dial(dest, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*200), grpc.MaxCallSendMsgSize(1024*1024*200)), grpc.WithTransportCredentials(insecure.NewCredentials()))
    stub := pb.NewMacroPodIngressClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    res, err := stub.EvalSummary(ctx, &ingress_eval_request)
    if err != nil {
        fmt.Println(err)
    }
    return res.GetReply()
}
