package function

import (
	"context"
	"fmt"
	"net/http"
	"bytes"
	"os"
	"encoding/json"
	"io/ioutil"
        "strconv"
        "math/rand"

	"time"
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

type RoomType struct {
        bookableRate float64
        totalRate float64
        totalRateInclusive float64
        code string
        currency string
        roomDescription string
}

type RatePlan struct {
        hotelId string
        code string
        inDate string
        outDate string
        roomType RoomType
}

type RatePlans []*RatePlan

type BodyGeo struct {
        Lat float64
        Lon float64
	WorkflowID string
        WorkflowDepth int
        WorkflowWidth int
}

type Request struct {
	HotelIds []string
	InDate string
	OutDate string
	WorkflowID string
        WorkflowDepth int
        WorkflowWidth int
}

// Nearby returns ids of nearby hotels ordered by ranking algo
func Nearby(req RequestBody) string {
        workflow_id := req.WorkflowID
        workflow_depth := req.WorkflowDepth
        workflow_width := req.WorkflowWidth
	// find nearby hotels
	fmt.Printf("in Search Nearby\n")

	fmt.Printf("nearby lat = %f\n", req.Lat)
	fmt.Printf("nearby lon = %f\n", req.Lon)

	requestURL := os.Getenv("HOTEL_GEO")
	payload := BodyGeo{Lat: req.Lat, Lon: req.Lon, WorkflowID: workflow_id, WorkflowDepth: workflow_depth + 1, WorkflowWidth: 0}
	body_g, err := json.Marshal(payload)
        req_url, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(body_g))
        req_url.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "1")
	nearby, err := client.Do(req_url)
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "2")
	if err != nil {
                fmt.Printf("nearby error: %v", err)
                return ""
        }

	// var ids []string
	nearbyBody, err := ioutil.ReadAll(nearby.Body)
	nearby_u := make([]string, 0)
	err = json.Unmarshal(nearbyBody, &nearby_u)
	for _, hid := range nearby_u {
		fmt.Printf("get Nearby hotelId = %s\n", hid)
		// ids = append(ids, hid)
	}

	// find rates for hotels
	r := Request{
		HotelIds: nearby_u,
		// HotelIds: []string{"2"},
		InDate:  req.InDate,
		OutDate: req.OutDate,
		WorkflowID: workflow_id,
                WorkflowDepth: workflow_depth + 1,
                WorkflowWidth: 0,
	}

	body_r, err := json.Marshal(r)

        requestURL2 := os.Getenv("HOTEL_RATE")
        req_url2, err := http.NewRequest(http.MethodPost, requestURL2, bytes.NewBuffer(body_r))
	req_url2.Header.Add("Content-Type", "application/json")
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "3")
        ratesRet, err := client.Do(req_url2)
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "4")
	if err != nil {
                fmt.Printf("rates error: %v", err)
                return ""
        }
	rates, err := ioutil.ReadAll(ratesRet.Body)
	// TODO(hw): add simple ranking algo to order hotel ids:
	// * geo distance
	// * price (best discount?)
	// * reviews

	// build the response
	res := make([]string, 0)
	var rate_p RatePlans
	json.Unmarshal(rates, &rate_p)
	for _, ratePlan := range rate_p {
		// fmt.Printf("get RatePlan HotelId = %s, Code = %s\n", ratePlan.HotelId, ratePlan.Code)
		res = append(res, ratePlan.hotelId)
	}
	ret, _ := json.Marshal(res)
        return string(ret)
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
        } else {
                body_u.WorkflowID = workflow_id
                body_u.WorkflowDepth = workflow_depth
                body_u.WorkflowWidth = workflow_width
        }
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "0")
	ret := Nearby(body_u)
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "5")
	fmt.Fprintf(res, ret) // echo to caller
}
