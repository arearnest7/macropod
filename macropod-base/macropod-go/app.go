package main

import (
    "os"
    "fmt"
    "time"
    "strconv"
    "net/http"
    "io/ioutil"
    "encoding/json"

    function "app/function"
)

func FunctionHandler(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    in := function.MacropodBody{}
    json.Unmarshal(body, &in)
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + in.WorkflowId + "," + strconv.Itoa(int(in.Depth)) + "," + strconv.Itoa(int(in.Width)) + "," + in.RequestType + "," + "request_start")
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + in.WorkflowId + "," + strconv.Itoa(int(in.Depth)) + "," + strconv.Itoa(int(in.Width)) + "," + in.RequestType + "," + "function_start")
    reply, _ := function.FunctionHandler(function.Context{Request: []byte(in.Data), WorkflowId: in.WorkflowId, Depth: int(in.Depth), Width: int(in.Width), RequestType: in.RequestType, InvokeType: "GRPC", IsJson: false})
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + in.WorkflowId + "," + strconv.Itoa(int(in.Depth)) + "," + strconv.Itoa(int(in.Width)) + "," + in.RequestType + "," + "request_end")
    fmt.Fprintf(res, reply)
}

func main() {
    http.HandleFunc("/", FunctionHandler)
    http.ListenAndServe(":" + os.Getenv("FUNC_PORT"), nil)
}
