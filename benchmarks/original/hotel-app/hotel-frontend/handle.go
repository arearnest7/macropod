package function

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"bytes"
	"encoding/json"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

type RequestBody struct {
        request string "json:\"request\""
	requestType string "json:\"requestType\""
        Lat float64 "json:\"Lat,omitempty\""
        Lon float64 "json:\"Lon,omitempty\""
        HotelId string "json:\"HotelId,omitempty\""
        HotelIds []string "json:\"HotelIds,omitempty\""
        RoomNumber int "json:\"RoomNumber,omitempty\""
        CustomerName string "json:\"CustomerName,omitempty\""
        Username string "json:\"Username,omitempty\""
        Password string "json:\"Password,omitempty\""
        Require string "json:\"Require,omitempty\""
        InDate string "json:\"InDate,omitempty\""
        OutDate string "json:\"OutDate,omitempty\""
}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	requestURL := ""
	body, err := ioutil.ReadAll(req.Body)
	var body_u *RequestBody
	json.Unmarshal(body, &body_u)
  	defer req.Body.Close()
	if body_u.request == "search" {
		requestURL = os.Getenv("HOTEL_SEARCH") + ":8080"
	} else if body_u.request == "recommend" {
                requestURL = os.Getenv("HOTEL_RECOMMEND") + ":8080"
	} else if body_u.request == "reserve" {
                requestURL = os.Getenv("HOTEL_RESERVE") + ":8080"
	} else if body_u.request == "user" {
                requestURL = os.Getenv("HOTEL_USER") + ":8080"
	}

	body_m, err := json.Marshal(body_u)
	req_url, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(body_m))
        if err != nil {
                log.Fatal(err)
        }
	req_url.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
        ret, err := client.Do(req_url)
        retBody, err := ioutil.ReadAll(ret.Body)
        ret_val, err := json.Marshal(retBody)
	fmt.Fprintf(res, string(ret_val)) // echo to caller
}

