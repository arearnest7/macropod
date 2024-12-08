package function

import (
    "fmt"
    "time"
    "strconv"
    "net/http"
    "encoding/json"
    "bytes"
    "io/ioutil"
)

type MacropodBody struct {
    Data string `json:"Data"`
    WorkflowId string `json:"WorkflowId"`
    Depth int32 `json:"Depth"`
    Width int32 `json:"Width"`
    RequestType string `json:"RequestType"`
    PvPath string `json:"PvPath"`
}

type Context struct {
    Request []byte
    WorkflowId string
    Depth int
    Width int
    RequestType string
    InvokeType string
    IsJson bool
}

func invoke(target string, ctx_in Context, data []byte, i int) (string) {
    request := MacropodBody{Data: string(data), WorkflowId: ctx_in.WorkflowId, Depth: int32(ctx_in.Depth + 1), Width: int32(i), RequestType: "gg"}
    body_m, _ := json.Marshal(request)
    req_url, err := http.NewRequest(http.MethodPost, "http://" + target, bytes.NewBuffer(body_m))
    if err != nil {
        fmt.Println(err)
    }
    req_url.Header.Add("Content-Type", "application/json")
    client := &http.Client{}
    ret, _ := client.Do(req_url)
    retBody, _ := ioutil.ReadAll(ret.Body)
    reply := string(retBody)

    return reply
}

func RPC(ctx_in Context, dest string, payloads [][]byte) ([]string) {
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.WorkflowId + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + "," + ctx_in.RequestType + "," + "rpc_start")
    tl := make(chan string)
    for i := 0; i < len(payloads); i++ {
        go func(n int) {
            tl <- invoke(dest, ctx_in, payloads[n], i)
        }(i)
    }
    results := make([]string, 0)
    for i := 0; i < len(payloads); i++ {
        res := <-tl
        results = append(results, res)
    }
    fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + ctx_in.WorkflowId + "," + strconv.Itoa(int(ctx_in.Depth)) + "," + strconv.Itoa(int(ctx_in.Width)) + "," + ctx_in.RequestType + "," + "rpc_end")
    return results
}
