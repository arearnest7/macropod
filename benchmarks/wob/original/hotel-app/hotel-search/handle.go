package function

import (
	"context"
	"fmt"
	"net/http"
	"bytes"
	"os"
	"encoding/json"
	"io/ioutil"

	"time"
        "github.com/redis/go-redis/v9"
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
}

type Request struct {
	HotelIds []string
	InDate string
	OutDate string
}

// Nearby returns ids of nearby hotels ordered by ranking algo
func Nearby(req RequestBody) string {
	// find nearby hotels
	fmt.Printf("in Search Nearby\n")

	fmt.Printf("nearby lat = %f\n", req.Lat)
	fmt.Printf("nearby lon = %f\n", req.Lon)

	requestURL := os.Getenv("HOTEL_GEO")
	payload := BodyGeo{Lat: req.Lat, Lon: req.Lon}
	body_g, err := json.Marshal(payload)
        req_url, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(body_g))
        req_url.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	nearby, err := client.Do(req_url)
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
	}

	body_r, err := json.Marshal(r)

        requestURL2 := os.Getenv("HOTEL_RATE")
        req_url2, err := http.NewRequest(http.MethodPost, requestURL2, bytes.NewBuffer(body_r))
	req_url2.Header.Add("Content-Type", "application/json")
        ratesRet, err := client.Do(req_url2)
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
	logging_name, logging := os.LookupEnv("LOGGING_NAME")
        redisClient := redis.NewClient(&redis.Options{})
        c := context.Background()
        if logging {
                logging_url := os.Getenv("LOGGING_URL")
                logging_password := os.Getenv("LOGGING_PASSWORD")
                redisClient = redis.NewClient(&redis.Options{
                        Addr: logging_url,
                        Password: logging_password,
                        DB: 0,
                })
        }
        if logging {
                redisClient.Append(c, logging_name, time.Now().String() + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n")
        }
        body, _ := ioutil.ReadAll(req.Body)
        body_u := RequestBody{}
        json.Unmarshal(body, &body_u)
        defer req.Body.Close()
	ret := Nearby(body_u)
	if logging {
                redisClient.Append(c, logging_name, time.Now().String() + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n")
        }
	fmt.Fprintf(res, ret) // echo to caller
}
