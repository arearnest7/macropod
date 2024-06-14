package function

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"bytes"
	"encoding/json"
	"io/ioutil"
        "strconv"
        "math/rand"

	"time"

	log "github.com/sirupsen/logrus"
)

type RequestBody struct {
        Request string `json:"Request"`
        RequestType string `json:"RequestType"`
        Lat float64 `json:"Lat"`
        Lon float64 `json:"Lon"`
        HotelId string `json:"HotelId"`
        HotelIds []string `json:"HotelIds"`
        RoomNumber int `json:"RoomNumber"`
        CustomerName string `json:"CustomerName"`
        Username string `json:"Username"`
        Password string `json:"Password"`
        Require string `json:"Require"`
        InDate string `json:"InDate"`
        OutDate string `json:"OutDate"`
	WorkflowID string `json:"WorkflowID"`
        WorkflowDepth int `json:"WorkflowDepth"`
        WorkflowWidth int `json:"WorkflowWidth"`
}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
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
                body_u.WorkflowDepth += 1
        } else {
                body_u.WorkflowID = workflow_id
                body_u.WorkflowDepth = workflow_depth + 1
                body_u.WorkflowWidth = workflow_width
        }
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "0")
	requestURL := ""
	if body_u.Request == "search" {
		requestURL = os.Getenv("HOTEL_SEARCH")
	} else if body_u.Request == "recommend" {
                requestURL = os.Getenv("HOTEL_RECOMMEND")
	} else if body_u.Request == "reserve" {
                requestURL = os.Getenv("HOTEL_RESERVE")
	} else if body_u.Request == "user" {
                requestURL = os.Getenv("HOTEL_USER")
	}

	body_m, err := json.Marshal(body_u)
	req_url, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(body_m))
        if err != nil {
                log.Fatal(err)
        }
	req_url.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "1")
        ret, err := client.Do(req_url)
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "2")
        retBody, err := ioutil.ReadAll(ret.Body)
        ret_val, err := json.Marshal(retBody)
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "3")
	fmt.Fprintf(res, string(ret_val)) // echo to caller
}

