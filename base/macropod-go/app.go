package main

import (
    pb "app/macropod_pb"
    structpb "google.golang.org/protobuf/types/known/structpb"
    "app/function"

    "os"
    "fmt"
    "strconv"
    "math/rand"
    "encoding/json"
    "io/ioutil"

    "net"
    "net/http"

    "google.golang.org/grpc"
    "golang.org/x/net/context"
)

type FunctionService struct {
    pb.UnimplementedMacroPodFunctionServer
}

var (

)

// purpose: serve the function invoke and return the results to the appropriate request handler
// input:
//  function_request: MacroPodRequest object containing the payload and all invocation related data
// output:
//  string: reply string
//  int: response code
//  error: any error result from the execution of the function; for now is nil
func Serve_Invoke(function_request *pb.MacroPodRequest) (string, int, error) {
    var reply string
    var code int
    var err error
    ctx_in := function.Context{Workflow: function_request.GetWorkflow(), Function: function_request.GetFunction(), Target: function_request.GetTarget(), Text: function_request.GetText(), JSON: function_request.GetJSON().AsMap(), Data: function_request.GetData(), WorkflowID: function_request.GetWorkflowID(), Depth: function_request.GetDepth(), Width: function_request.GetWidth()}
    function.Timestamp(ctx_in, "", "request")
    reply, code = function.FunctionHandler(ctx_in)
    function.Timestamp(ctx_in, "", "request_end")
    return reply, code, err
}

// purpose: gRPC handler for Invoke
// input:
//  ctx: context object utilized in every gRPC request handler
//  req: MacroPodRequest object defined by protobuf
// output:
//  MacroPodReply: reply protobuf struct
//  error: any error printout from invocation
func (s *FunctionService) Invoke(ctx context.Context, in *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
    results, code, err := Serve_Invoke(in)
    c := int32(code)
    res := pb.MacroPodReply{Reply: &results, Code: &c}
    return &res, err
}

// purpose: http handler for Invoke
// input:
//  res: writer to provide HTTP reply data
//  req: GET/POST object
func HTTP_Invoke(res http.ResponseWriter, req *http.Request) {
    workflow := os.Getenv("WORKFLOW")
    function := os.Getenv("FUNCTION")
    workflow_id := strconv.Itoa(rand.Intn(100000))
    dw := int32(0)
    target := req.Header.Get("Host")
    body, _ := ioutil.ReadAll(req.Body)
    in := pb.MacroPodRequest{Workflow: &workflow, Function: &function, WorkflowID: &workflow_id, Depth: &dw, Width: &dw, Target: &target}
    switch content := req.Header.Get("Content-Type"); content {
        case "text/plain":
            t := string(body)
            in.Text = &t
        case "application/json":
            j_i := make(map[string]interface{})
            json.Unmarshal(body, &j_i)
            j, _ := structpb.NewStruct(j_i)
            in.JSON = j
        case "application/octet-stream":
            in.Data = body
        default:
            t := string(body)
            in.Text = &t
    }
    results, _, _ := Serve_Invoke(&in)
    fmt.Fprintf(res, results)
}

// input:
//  SERVICE_PORT: string defining gRPC service binding; default 5000
//  HTTP_PORT: string defining HTTP service binding; default 6000
//  WORKFLOW: string name of workflow the function belongs to
//  FUNCTION: string name of the function
//  INGRESS: string representing the IP:port of the ingress such that gateway communication can be used if enabled
//  LOGGER: string IP:port of the logger if logging is enabled; empty if logging is disabled
//  COMM_TYPE: string name of the communication type (refer to the deployer documentation)
func main() {
    service_port := os.Getenv("SERVICE_PORT")
    if service_port == "" {
        service_port = "5000"
    }
    http_port := os.Getenv("HTTP_PORT")
    if http_port == "" {
        http_port = "6000"
    }
    l, _ := net.Listen("tcp", ":" + service_port)
    s := grpc.NewServer(grpc.MaxSendMsgSize(1024*1024*200), grpc.MaxRecvMsgSize(1024*1024*200))
    pb.RegisterMacroPodFunctionServer(s, &FunctionService{})

    go s.Serve(l)

    h := http.NewServeMux()
    h.HandleFunc("/", HTTP_Invoke)
    http.ListenAndServe(":" + http_port, h)
}
