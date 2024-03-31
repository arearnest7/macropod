package function

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"bytes"
	"encoding/json"
	"io/ioutil"

	"time"
        "github.com/redis/go-redis/v9"

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
}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	logging_name, logging := os.LookupEnv("LOGGING_NAME")
        redisClient := redis.NewClient(&redis.Options{})
        c := context.Background()
        if logging {
                logging_ip := os.Getenv("LOGGING_IP")
                logging_password := os.Getenv("LOGGING_PASSWORD")
                redisClient = redis.NewClient(&redis.Options{
                        Addr: logging_ip,
                        Password: logging_password,
                        DB: 0,
                })
        }
        if logging {
                redisClient.Append(c, logging_name, time.Now().String() + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n")
        }
	requestURL := ""
	body, err := ioutil.ReadAll(req.Body)
	body_u := RequestBody{}
	json.Unmarshal(body, &body_u)
  	defer req.Body.Close()
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
        ret, err := client.Do(req_url)
        retBody, err := ioutil.ReadAll(ret.Body)
        ret_val, err := json.Marshal(retBody)
	if logging {
                redisClient.Append(c, logging_name, time.Now().String() + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n")
        }
	fmt.Fprintf(res, string(ret_val)) // echo to caller
}

