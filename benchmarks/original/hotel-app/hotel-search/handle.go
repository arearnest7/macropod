package function

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"net"

	log "github.com/sirupsen/logrus"

	"time"

	"golang.org/x/net/context"
)

type RoomType struct {
        bookableRate double
        totalRate double
        totalRateInclusive double
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

// Nearby returns ids of nearby hotels ordered by ranking algo
func Nearby(var req) (string, error) {
	// find nearby hotels
	fmt.Printf("in Search Nearby\n")

	fmt.Printf("nearby lat = %f\n", req.Lat)
	fmt.Printf("nearby lon = %f\n", req.Lon)

	contents, err := ioutil.ReadFile("/etc/secret-volume/hotel-geo")
        if err != nil {
                log.Fatal(err)
        }
        hotel-geo := string(contents)

	requestURL := hotel-geo + ":80"
        req_url, err := http.NewRequest(http.MethodPost, requestURL, {Lat: req.Lat, Lon: req.Lon})
        nearby, err := http.DefaultClient.Do(req_url)
	if err != nil {
                fmt.Printf("nearby error: %v", err)
                return nil, err
        }

	// var ids []string
	for _, hid := range json.Unmarshal(nearby) {
		fmt.Printf("get Nearby hotelId = %s\n", hid)
		// ids = append(ids, hid)
	}

	// find rates for hotels
	r := rate.Request{
		HotelIds: json.Unmarshal(nearby),
		// HotelIds: []string{"2"},
		InDate:  req.InDate,
		OutDate: req.OutDate,
	}

	contents, err := ioutil.ReadFile("/etc/secret-volume/hotel-rate")
        if err != nil {
                log.Fatal(err)
        }
        hotel-rate := string(contents)

        requestURL := hotel-rate + ":80"
        req_url, err := http.NewRequest(http.MethodPost, requestURL, &r)
        rates, err := http.DefaultClient.Do(req_url)
	if err != nil {
                fmt.Printf("rates error: %v", err)
                return nil, err
        }
	// TODO(hw): add simple ranking algo to order hotel ids:
	// * geo distance
	// * price (best discount?)
	// * reviews

	// build the response
	res := make([]string, 0)
	rate_p : make(RatePlans, 1)
	json.Unmarshal(rates, &rate_p)
	for _, ratePlan := range rate_p {
		// fmt.Printf("get RatePlan HotelId = %s, Code = %s\n", ratePlan.HotelId, ratePlan.Code)
		res = append(res, ratePlan.HotelId)
	}
	return json.Marshal(res), nil
}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
        fmt.Fprintf(res, Nearby(req.json)) // echo to caller
}
