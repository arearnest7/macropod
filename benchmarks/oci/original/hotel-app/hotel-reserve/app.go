package main

import (
    "net/http"
    "time"
    "fmt"
    "os"
    "strconv"
    "math/rand"
    "encoding/json"
    "io/ioutil"
)

func function(res http.ResponseWriter, req *http.Request) {
    workflow_id := strconv.Itoa(rand.Intn(10000000))
    workflow_depth := 0
    workflow_width := 0
    body, _ := ioutil.ReadAll(req.Body)
    body_u := RequestBody{}
    json.Unmarshal(body, &body_u)
    defer req.Body.Close()
    if body_u.WorkflowID != "" {
        workflow_id = body_u.WorkflowID
        workflow_depth = body_u.WorkflowDepth
        workflow_width = body_u.WorkflowWidth
    } else {
        body_u.WorkflowID = workflow_id
        body_u.WorkflowDepth = workflow_depth
        body_u.WorkflowWidth = workflow_width
    }
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "0")
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "1")
    FunctionHandler(res, body_u, "application/json", true)
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "2")
}

func main() {
    http.HandleFunc("/", function)
    http.ListenAndServe(":" + os.Getenv("PORT"), nil)
}
