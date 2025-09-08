package main

import (
    pb "app/macropod_pb"

    "os"
    "fmt"
    "strconv"
    "encoding/json"
    "io/ioutil"
    "sync"
    "time"

    "net"
    "net/http"

    "google.golang.org/grpc"
    "golang.org/x/net/context"
)

type LoggerService struct {
    pb.UnimplementedMacroPodLoggerServer
}

var (
    timestamp_dir  string                                       // string name of the timestamp directory
    timestamp_lock = make(map[string]map[string]*sync.Mutex)    // mutex list for every timestamp file to prevent parallel read/writes
    error_dir      string                                       // string name of the error directory
    error_lock     = make(map[string]map[string]*sync.Mutex)    // mutex list for every error file to prevent parallel read/writes
    print_dir      string                                       // string name of the print directory
    print_lock     = make(map[string]map[string]*sync.Mutex)    // mutex list for every print file to prevent parallel read/writes
    dataLock       sync.Mutex                                   // general purpose mutex lock for all operations modifying global variables (preventing system panics and divergent data commits)
)

// purpose: append to timestamp file for the workflow:function combination the sent message and invocation data
// input:
//  request: MacroPodRequest that contains all message data and invocation data from function
func Serve_Timestamp(request *pb.MacroPodRequest) {
    dataLock.Lock()
    _, exists := timestamp_lock[request.GetWorkflow()]
    if !exists {
        timestamp_lock[request.GetWorkflow()] = make(map[string]*sync.Mutex)
    }
    _, exists = timestamp_lock[request.GetWorkflow()][request.GetFunction()]
    if !exists {
        m := sync.Mutex{}
        timestamp_lock[request.GetWorkflow()][request.GetFunction()] = &m
    }
    dataLock.Unlock()
    timestamp_lock[request.GetWorkflow()][request.GetFunction()].Lock()
    if _, err := os.Stat(timestamp_dir + request.GetWorkflow()); err != nil {
        if os.IsNotExist(err) {
            os.Mkdir(timestamp_dir + request.GetWorkflow(), 0777)
        } else {
            fmt.Println(err)
        }
    }
    f, _ := os.OpenFile(timestamp_dir + request.GetWorkflow() + "/" + request.GetFunction(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
    f.WriteString(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + request.GetWorkflowID() + "," + strconv.Itoa(int(request.GetDepth())) + "," + strconv.Itoa(int(request.GetWidth())) + "," + request.GetTarget() + "," + request.GetText() + "\n")
    f.Close()
    timestamp_lock[request.GetWorkflow()][request.GetFunction()].Unlock()
}

// purpose: append to error file for the workflow:function combination the sent message
// input:
//  request: MacroPodRequest that contains all message data and invocation data from function
func Serve_Error(request *pb.MacroPodRequest) {
    dataLock.Lock()
    _, exists := error_lock[request.GetWorkflow()]
    if !exists {
        error_lock[request.GetWorkflow()] = make(map[string]*sync.Mutex)
    }
    _, exists = error_lock[request.GetWorkflow()][request.GetFunction()]
    if !exists {
        m := sync.Mutex{}
        error_lock[request.GetWorkflow()][request.GetFunction()] = &m
    }
    dataLock.Unlock()
    error_lock[request.GetWorkflow()][request.GetFunction()].Lock()
    if _, err := os.Stat(error_dir + request.GetWorkflow()); err != nil {
        if os.IsNotExist(err) {
            os.Mkdir(error_dir + request.GetWorkflow(), 0777)
        } else {
            fmt.Println(err)
        }
    }
    f, _ := os.OpenFile(error_dir + request.GetWorkflow() + "/" + request.GetFunction(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
    f.WriteString(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + ": " + request.GetText() + "\n")
    f.Close()
    error_lock[request.GetWorkflow()][request.GetFunction()].Unlock()
}

// purpose: append to print file for the workflow:function combination the sent message
// input:
//  request: MacroPodRequest that contains all message data and invocation data from function
func Serve_Print(request *pb.MacroPodRequest) {
    dataLock.Lock()
    _, exists := print_lock[request.GetWorkflow()]
    if !exists {
        print_lock[request.GetWorkflow()] = make(map[string]*sync.Mutex)
    }
    _, exists = print_lock[request.GetWorkflow()][request.GetFunction()]
    if !exists {
        m := sync.Mutex{}
        print_lock[request.GetWorkflow()][request.GetFunction()] = &m
    }
    dataLock.Unlock()
    print_lock[request.GetWorkflow()][request.GetFunction()].Lock()
    if _, err := os.Stat(print_dir + request.GetWorkflow()); err != nil {
        if os.IsNotExist(err) {
            os.Mkdir(print_dir + request.GetWorkflow(), 0777)
        } else {
            fmt.Println(err)
        }
    }
    f, _ := os.OpenFile(print_dir + request.GetWorkflow() + "/" + request.GetFunction(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
    f.WriteString(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + ": " + request.GetText() + "\n")
    f.Close()
    print_lock[request.GetWorkflow()][request.GetFunction()].Unlock()
}

// purpose: retrieve all timestamp logs from the function and workflow specified
// input:
//  request: MacroPodRequest that contains the workflow and function information for retrieval
// output:
//  string: CSV file in string form
func Serve_GetTimestamp(request *pb.MacroPodRequest) (string) {
    dataLock.Lock()
    _, exists := timestamp_lock[request.GetWorkflow()]
    if !exists {
        timestamp_lock[request.GetWorkflow()] = make(map[string]*sync.Mutex)
    }
    _, exists = timestamp_lock[request.GetWorkflow()][request.GetFunction()]
    if !exists {
        m := sync.Mutex{}
        timestamp_lock[request.GetWorkflow()][request.GetFunction()] = &m
    }
    dataLock.Unlock()
    timestamp_lock[request.GetWorkflow()][request.GetFunction()].Lock()
    if _, err := os.Stat(timestamp_dir + request.GetWorkflow()); err != nil {
        if os.IsNotExist(err) {
            os.Mkdir(timestamp_dir + request.GetWorkflow(), 0777)
        } else {
            fmt.Println(err)
        }
    }
    b, _ := os.ReadFile(timestamp_dir + request.GetWorkflow() + "/" + request.GetFunction())
    timestamp_lock[request.GetWorkflow()][request.GetFunction()].Unlock()
    return string(b)
}

// purpose: retrieve all error logs from the function and workflow specified
// input:
//  request: MacroPodRequest that contains the workflow and function information for retrieval
// output:
//  string: CSV file in string form
func Serve_GetError(request *pb.MacroPodRequest) (string) {
    dataLock.Lock()
    _, exists := error_lock[request.GetWorkflow()]
    if !exists {
        error_lock[request.GetWorkflow()] = make(map[string]*sync.Mutex)
    }
    _, exists = error_lock[request.GetWorkflow()][request.GetFunction()]
    if !exists {
        m := sync.Mutex{}
        error_lock[request.GetWorkflow()][request.GetFunction()] = &m
    }
    dataLock.Unlock()
    error_lock[request.GetWorkflow()][request.GetFunction()].Lock()
    if _, err := os.Stat(error_dir + request.GetWorkflow()); err != nil {
        if os.IsNotExist(err) {
            os.Mkdir(error_dir + request.GetWorkflow(), 0777)
        } else {
            fmt.Println(err)
        }
    }
    b, _ := os.ReadFile(error_dir + request.GetWorkflow() + "/" + request.GetFunction())
    error_lock[request.GetWorkflow()][request.GetFunction()].Unlock()
    return string(b)
}

// purpose: retrieve all print logs from the function and workflow specified
// input:
//  request: MacroPodRequest that contains the workflow and function information for retrieval
// output:
//  string: CSV file in string form
func Serve_GetPrint(request *pb.MacroPodRequest) (string) {
    dataLock.Lock()
    _, exists := print_lock[request.GetWorkflow()]
    if !exists {
        print_lock[request.GetWorkflow()] = make(map[string]*sync.Mutex)
    }
    _, exists = print_lock[request.GetWorkflow()][request.GetFunction()]
    if !exists {
        m := sync.Mutex{}
        print_lock[request.GetWorkflow()][request.GetFunction()] = &m
    }
    dataLock.Unlock()
    print_lock[request.GetWorkflow()][request.GetFunction()].Lock()
    if _, err := os.Stat(print_dir + request.GetWorkflow()); err != nil {
        if os.IsNotExist(err) {
            os.Mkdir(print_dir + request.GetWorkflow(), 0777)
        } else {
            fmt.Println(err)
        }
    }
    b, _ := os.ReadFile(print_dir + request.GetWorkflow() + "/" + request.GetFunction())
    print_lock[request.GetWorkflow()][request.GetFunction()].Unlock()
    return string(b)
}

// purpose: gRPC handler for Timestamp
// input:
//  ctx: context object utilized in every gRPC request handler
//  req: MacroPodRequest object defined by protobuf
// output:
//  MacroPodReply: reply protobuf struct
//  error: always will be nil
func (s *LoggerService) Timestamp(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
    Serve_Timestamp(req)
    results := pb.MacroPodReply{}
    return &results, nil
}

// purpose: gRPC handler for Error
// input:
//  ctx: context object utilized in every gRPC request handler
//  req: MacroPodRequest object defined by protobuf
// output:
//  MacroPodReply: reply protobuf struct
//  error: always will be nil
func (s *LoggerService) Error(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
    Serve_Error(req)
    results := pb.MacroPodReply{}
    return &results, nil
}

// purpose: gRPC handler for Print
// input:
//  ctx: context object utilized in every gRPC request handler
//  req: MacroPodRequest object defined by protobuf
// output:
//  MacroPodReply: reply protobuf struct
//  error: always will be nil
func (s *LoggerService) Print(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
    Serve_Print(req)
    results := pb.MacroPodReply{}
    return &results, nil
}

// purpose: gRPC handler for GetTimestamp
// input:
//  ctx: context object utilized in every gRPC request handler
//  req: MacroPodRequest object defined by protobuf
// output:
//  MacroPodReply: reply protobuf struct
//  error: always will be nil
func (s *LoggerService) GetTimestamp(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
    t := Serve_GetTimestamp(req)
    results := pb.MacroPodReply{Reply: &t}
    return &results, nil
}

// purpose: gRPC handler for GetError
// input:
//  ctx: context object utilized in every gRPC request handler
//  req: MacroPodRequest object defined by protobuf
// output:
//  MacroPodReply: reply protobuf struct
//  error: always will be nil
func (s *LoggerService) GetError(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
    e := Serve_GetError(req)
    results := pb.MacroPodReply{Reply: &e}
    return &results, nil
}

// purpose: gRPC handler for GetPrint
// input:
//  ctx: context object utilized in every gRPC request handler
//  req: MacroPodRequest object defined by protobuf
// output:
//  MacroPodReply: reply protobuf struct
//  error: always will be nil
func (s *LoggerService) GetPrint(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
    p := Serve_GetPrint(req)
    results := pb.MacroPodReply{Reply: &p}
    return &results, nil
}

// purpose: prints out a help message for quick reference on how to operate service
// input:
//  res: writer to provide HTTP reply data
//  req: GET/POST object
func HTTP_Help(res http.ResponseWriter, req *http.Request) {
    help_print := "usage:\n"
    help_print += "\t/timestamp/set/{workflow}/{function}\n"
    help_print += "\t\tpurpose: append a timestamp message to the log for specified function of specified workflow\n"
    help_print += "\t\tinput: JSON object of MacroPodRequest\n"
    help_print += "\t/error/set/{workflow}/{function}\n"
    help_print += "\t\tpurpose: append a error message to the log for specified function of specified workflow\n"
    help_print += "\t\tinput: JSON object of MacroPodRequest\n"
    help_print += "\t/print/set/{workflow}/{function}\n"
    help_print += "\t\tpurpose: append a print message to the log for specified function of specified workflow\n"
    help_print += "\t\tinput: JSON object of MacroPodRequest\n"
    help_print += "\t/timestamp/get/{workflow}/{function}\n"
    help_print += "\t\tpurpose: retrieve all timestamp logs for specified function of specified workflow\n"
    help_print += "\t\toutput: CSV file\n"
    help_print += "\t/error/get/{workflow}/{function}\n"
    help_print += "\t\tpurpose: retrieve all error logs for specified function of specified workflow\n"
    help_print += "\t\toutput: CSV file\n"
    help_print += "\t/print/get/{workflow}/{function}\n"
    help_print += "\t\tpurpose: retrieve all print logs for specified function of specified workflow\n"
    help_print += "\t\toutput: CSV file\n"
    fmt.Fprint(res, help_print)
}

// purpose: http handler for Timestamp
// input:
//  res: writer to provide HTTP reply data
//  req: GET/POST object
func HTTP_Timestamp(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    request := pb.MacroPodRequest{}
    json.Unmarshal(body, &request)
    Serve_Timestamp(&request)
}

// purpose: http handler for Error
// input:
//  res: writer to provide HTTP reply data
//  req: GET/POST object
func HTTP_Error(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    request := pb.MacroPodRequest{}
    json.Unmarshal(body, &request)
    Serve_Error(&request)
}

// purpose: http handler for Print
// input:
//  res: writer to provide HTTP reply data
//  req: GET/POST object
func HTTP_Print(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    request := pb.MacroPodRequest{}
    json.Unmarshal(body, &request)
    Serve_Print(&request)
}

// purpose: http handler for GetTimestamp
// input:
//  res: writer to provide HTTP reply data
//  req: GET/POST object
func HTTP_GetTimestamp(res http.ResponseWriter, req *http.Request) {
    workflow := req.PathValue("workflow")
    function := req.PathValue("function")
    request := pb.MacroPodRequest{Workflow: &workflow, Function: &function}
    results := Serve_GetTimestamp(&request)
    fmt.Fprint(res, results)
}

// purpose: http handler for GetError
// input:
//  res: writer to provide HTTP reply data
//  req: GET/POST object
func HTTP_GetError(res http.ResponseWriter, req *http.Request) {
    workflow := req.PathValue("workflow")
    function := req.PathValue("function")
    request := pb.MacroPodRequest{Workflow: &workflow, Function: &function}
    results := Serve_GetError(&request)
    fmt.Fprint(res, results)
}

// purpose: http handler for GetPrint
// input:
//  res: writer to provide HTTP reply data
//  req: GET/POST object
func HTTP_GetPrint(res http.ResponseWriter, req *http.Request) {
    workflow := req.PathValue("workflow")
    function := req.PathValue("function")
    request := pb.MacroPodRequest{Workflow: &workflow, Function: &function}
    results := Serve_GetPrint(&request)
    fmt.Fprint(res, results)
}

// input:
//  SERVICE_PORT: string defining gRPC service binding; default 8000
//  HTTP_PORT: string defining HTTP service binding; default 9000
//  TIMESTAMP_DIR: directory for storing all timestamp data; default /app/timestamp/
//  ERROR_DIR: directory for storing all error data; default /app/error/
//  PRINT_DIR: directory for storing all print data; default /app/print/
func main() {
    service_port := os.Getenv("SERVICE_PORT")
    if service_port == "" {
        service_port = "8000"
    }
    http_port := os.Getenv("HTTP_PORT")
    if http_port == "" {
        http_port = "9000"
    }
    timestamp_dir = os.Getenv("TIMESTAMP_DIR")
    if timestamp_dir == "" {
        timestamp_dir = "/app/timestamp/"
    }
    error_dir = os.Getenv("ERROR_DIR")
    if error_dir == "" {
        error_dir = "/app/error/"
    }
    print_dir = os.Getenv("PRINT_DIR")
    if print_dir == "" {
        print_dir = "/app/print/"
    }
    if _, err := os.Stat(timestamp_dir); err != nil {
        if os.IsNotExist(err) {
            os.Mkdir(timestamp_dir, 0777)
        } else {
            fmt.Println(err)
        }
    }
    if _, err := os.Stat(error_dir); err != nil {
        if os.IsNotExist(err) {
            os.Mkdir(error_dir, 0777)
        } else {
            fmt.Println(err)
        }
    }
    if _, err := os.Stat(print_dir); err != nil {
        if os.IsNotExist(err) {
            os.Mkdir(print_dir, 0777)
        } else {
            fmt.Println(err)
        }
    }
    l, _ := net.Listen("tcp", ":" + service_port)
    s := grpc.NewServer(grpc.MaxSendMsgSize(1024*1024*200), grpc.MaxRecvMsgSize(1024*1024*200))
    pb.RegisterMacroPodLoggerServer(s, &LoggerService{})

    go s.Serve(l)

    h := http.NewServeMux()
    h.HandleFunc("/", HTTP_Help)
    h.HandleFunc("/timestamp/set/{workflow}/{function}", HTTP_Timestamp)
    h.HandleFunc("/error/set/{workflow}/{function}", HTTP_Error)
    h.HandleFunc("/print/set/{workflow}/{function}", HTTP_Print)
    h.HandleFunc("/timestamp/get/{workflow}/{function}", HTTP_GetTimestamp)
    h.HandleFunc("/error/get/{workflow}/{function}", HTTP_GetError)
    h.HandleFunc("/print/get/{workflow}/{function}", HTTP_GetPrint)
    http.ListenAndServe(":" + http_port, h)
}
