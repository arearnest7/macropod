package main

import (
    pb "app/macropod_pb"
    "app/function"

    "os"
    "encoding/json"
    "fmt"
    "strconv"
    "time"
    "io/ioutil"

    "github.com/soheilhy/cmux"
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

func Function_Invoke(function_request pb.FunctionRequest) (string, int, error) {
    var reply string
    var code int
    var err error
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + function_request.GetInvocationID() + "," + strconv.Itoa(int(function_request.GetDepth())) + "," + strconv.Itoa(int(function_request.GetWidth())) + "," + "request_start")
    defer func() {
        if err := recover(); err != nil {}
    }()
    reply, code = function.FunctionHandler(function.Context{Function: function_request.GetFunction(), Data: function_request.GetData(), Text: function_request.GetText(), InvocationID: function_request.GetInvocationID(), Depth: function_request.GetDepth(), Width: function_request.GetWidth()})
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + function_request.GetInvocationID() + "," + strconv.Itoa(int(function_request.GetDepth())) + "," + strconv.Itoa(int(function_request.GetWidth())) + "," + "request_end")
    return reply, code, err
}

func (s *FunctionService) Invoke(ctx context.Context, in *pb.FunctionRequest) (*pb.MacroPodReply, error) {
    results, code, err := Function_Invoke(*in)
    c := int32(code)
    res := pb.MacroPodReply{Reply: &results, Code: &c}
    return &res, err
}

func Serve_Function_Invoke(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    function_request := pb.FunctionRequest{}
    json.Unmarshal(body, &function_request)
    defer req.Body.Close()
    results, _, err := Function_Invoke(function_request)
    if err != nil {
        fmt.Println(res, err)
    }
    fmt.Fprint(res, results)
}

func main() {
    service_port := os.Getenv("SERVICE_PORT")
    if service_port == "" {
        service_port = "8081"
    }
    l, err := net.Listen("tcp", ":" + service_port)
    if err != nil {
        fmt.Println("listener not operational")
    }
    m := cmux.New(l)
    l_g := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
    l_h := m.Match(cmux.HTTP1Fast())

    s_g := grpc.NewServer()
    m_h := http.NewServeMux()
    pb.RegisterMacroPodFunctionServer(s_g, &FunctionService{})
    m_h.HandleFunc("/", Serve_Function_Invoke)

    go s_g.Serve(l_g)
    s_h := &http.Server{Handler: m_h}
    go s_h.Serve(l_h)
}
