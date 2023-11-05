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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	geo "github.com/vhive-serverless/vSwarm-proto/proto/hotel_reserv/geo"
	rate "github.com/vhive-serverless/vSwarm-proto/proto/hotel_reserv/rate"

	pb "github.com/vhive-serverless/vSwarm-proto/proto/hotel_reserv/search"

	tracing "github.com/vhive-serverless/vSwarm/utils/tracing/go"
)

// Nearby returns ids of nearby hotels ordered by ranking algo
func Nearby(var req) (*pb.SearchResult, error) {
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
	for _, hid := range nearby.HotelIds {
		fmt.Printf("get Nearby hotelId = %s\n", hid)
		// ids = append(ids, hid)
	}

	// find rates for hotels
	r := rate.Request{
		HotelIds: nearby.HotelIds,
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
	res := new(pb.SearchResult)
	for _, ratePlan := range rates.RatePlans {
		// fmt.Printf("get RatePlan HotelId = %s, Code = %s\n", ratePlan.HotelId, ratePlan.Code)
		res.HotelIds = append(res.HotelIds, ratePlan.HotelId)
	}
	return res, nil
}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
        fmt.Fprintf(res, Nearby(req.json)) // echo to caller
}
