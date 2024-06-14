package main

import (
    "net/http"
    "os"
    "context"
    "fmt"
    "net"
    "time"
    "math/rand"
    "strconv"
    "io/ioutil"
    "github.com/go-mmap/mmap"
    "google.golang.org/grpc"

    pb "app/wf_pb"
    function "app/function"
)

type server struct {
    pb.UnimplementedGRPCFunctionServer
}

func HTTPFunctionHandler(res http.ResponseWriter, req *http.Request) {
    workflow_id := strconv.Itoa(rand.Intn(10000000))
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + "0" + "," + "0" + "," + "HTTP" + "," + "0")
    body, _ := ioutil.ReadAll(req.Body)
    request_type := "gg"
    _, pv := os.LookupEnv("APP_PV")
    if pv {
        request_type = "gm"
    }
    if req.Header.Get("Content-Type") == "application/json" {
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + "0" + "," + "0" + "," + "HTTP" + "," + "1")
        reply, _ := function.FunctionHandler(function.Context{Request: body, WorkflowId: workflow_id, Depth: 0, Width: 0, RequestType: request_type, InvokeType: "HTTP", IsJson: true})
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + "0" + "," + "0" + "," + "HTTP" + "," + "2")
        fmt.Fprintf(res, reply)
    } else {
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + "0" + "," + "0" + "," + "HTTP" + "," + "1")
        reply, _ := function.FunctionHandler(function.Context{Request: body, WorkflowId: workflow_id, Depth: 0, Width: 0, RequestType: request_type, InvokeType: "HTTP", IsJson: false})
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + "0" + "," + "0" + "," + "HTTP" + "," + "2")
        fmt.Fprintf(res, reply)
    }
}

func (s *server) GRPCFunctionHandler(ctx context.Context, in *pb.RequestBody) (*pb.ResponseBody, error) {
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + in.GetWorkflowId() + "," + strconv.Itoa(int(in.GetDepth())) + "," + strconv.Itoa(int(in.GetWidth())) + "," + in.GetRequestType() + "," + "0")
    app_pv, _ := os.LookupEnv("APP_PV")
    res := pb.ResponseBody{Code: int32(500)}
    if in.GetRequestType() == "" || in.GetRequestType() == "gg" {
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + in.GetWorkflowId() + "," + strconv.Itoa(int(in.GetDepth())) + "," + strconv.Itoa(int(in.GetWidth())) + "," + in.GetRequestType() + "," + "1")
        reply, code := function.FunctionHandler(function.Context{Request: in.GetData(), WorkflowId: in.GetWorkflowId(), Depth: int(in.GetDepth()), Width: int(in.GetWidth()), RequestType: in.GetRequestType(), InvokeType: "GRPC", IsJson: false})
        res.Reply = &reply
        res.Code = int32(code)
    } else if in.GetRequestType() == "mg" {
        f, _ := mmap.Open(app_pv + "/" + in.GetPvPath())
        req := make([]byte, f.Len())
        f.ReadAt(req, 0)
        f.Close()
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + in.GetWorkflowId() + "," + strconv.Itoa(int(in.GetDepth())) + "," + strconv.Itoa(int(in.GetWidth())) + "," + in.GetRequestType() + "," + "1")
        reply, code := function.FunctionHandler(function.Context{Request: req, WorkflowId: in.GetWorkflowId(), Depth: int(in.GetDepth()), Width: int(in.GetWidth()), RequestType: in.GetRequestType(), InvokeType: "GRPC", IsJson: false})
        res.Reply = &reply
        res.Code = int32(code)
    } else if in.GetRequestType() == "gm" {
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + in.GetWorkflowId() + "," + strconv.Itoa(int(in.GetDepth())) + "," + strconv.Itoa(int(in.GetWidth())) + "," + in.GetRequestType() + "," + "1")
        payload, code := function.FunctionHandler(function.Context{Request: in.GetData(), WorkflowId: in.GetWorkflowId(), Depth: int(in.GetDepth()), Width: int(in.GetWidth()), RequestType: in.GetRequestType(), InvokeType: "GRPC", IsJson: false})
        pv_path := in.GetWorkflowId() + "_" + strconv.Itoa(int(in.GetDepth())) + "_" + strconv.Itoa(int(in.GetWidth())) + "_" + strconv.Itoa(rand.Intn(10000000))
        os.WriteFile(app_pv + "/" + pv_path, []byte(payload), 777)
        reply := ""
        res.Reply = &reply
        res.Code = int32(code)
        res.PvPath = &pv_path
    } else if in.GetRequestType() == "mm" {
        f, _ := mmap.Open(app_pv + "/" + in.GetPvPath())
        req := make([]byte, f.Len())
        f.ReadAt(req, 0)
        f.Close()
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + in.GetWorkflowId() + "," + strconv.Itoa(int(in.GetDepth())) + "," + strconv.Itoa(int(in.GetWidth())) + "," + in.GetRequestType() + "," + "1")
        payload, code := function.FunctionHandler(function.Context{Request: req, WorkflowId: in.GetWorkflowId(), Depth: int(in.GetDepth()), Width: int(in.GetWidth()), RequestType: in.GetRequestType(), InvokeType: "GRPC", IsJson: false})
        pv_path := in.GetWorkflowId() + "_" + strconv.Itoa(int(in.GetDepth())) + "_" + strconv.Itoa(int(in.GetWidth())) + "_" + strconv.Itoa(rand.Intn(10000000))
        os.WriteFile(app_pv + "/" + pv_path, []byte(payload), 777)
        reply := ""
        res.Reply = &reply
        res.Code = int32(code)
        res.PvPath = &pv_path
    }
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + in.GetWorkflowId() + "," + strconv.Itoa(int(in.GetDepth())) + "," + strconv.Itoa(int(in.GetWidth())) + "," + in.GetRequestType() + "," + "2")
    return &res, nil
}

func main() {
    service_type, typed := os.LookupEnv("SERVICE_TYPE")
    if !typed || service_type == "HTTP" {
        http.HandleFunc("/", HTTPFunctionHandler)
        http.ListenAndServe(":" + os.Getenv("FUNC_PORT"), nil)
    } else if service_type == "GRPC" {
        port, _ := strconv.Atoi(os.Getenv("FUNC_PORT"))
        l, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))
        s := grpc.NewServer(grpc.MaxSendMsgSize(1024*1024*200), grpc.MaxRecvMsgSize(1024*1024*200))
        pb.RegisterGRPCFunctionServer(s, &server{})
        s.Serve(l)
    }
}
