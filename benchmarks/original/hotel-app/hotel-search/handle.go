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
)

type RequestBody struct {
        request string "json:\"request\""
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
	Lat string
	Lon string
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

	requestURL := os.Getenv("HOTEL_GEO") + ":8080"
	payload := BodyGeo{Lat: strconv.FormatFloat(req.Lat, 'f', -1, 64), Lon: strconv.FormatFloat(req.Lon, 'f', -1, 64)}
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
	var nearby_u []string
	err = json.Unmarshal(nearbyBody, nearby_u)
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

        requestURL2 := os.Getenv("HOTEL_RATE") + ":8080"
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
        body, _ := ioutil.ReadAll(req.Body)
        var body_u *RequestBody
        json.Unmarshal(body, &body_u)
        defer req.Body.Close()
	fmt.Fprintf(res, Nearby(*body_u)) // echo to caller
}
