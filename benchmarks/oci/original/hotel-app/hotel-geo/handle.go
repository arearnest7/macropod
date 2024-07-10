package main

import (
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
	"encoding/json"
	"time"
        "strconv"

	log "github.com/sirupsen/logrus"

	"github.com/hailocab/go-geoindex"
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

const (
	maxSearchRadius  = 10
	maxSearchResults = 5
)

type Point struct {
	Pid  string  `bson:"hotelid"`
	Plat float64 `bson:"lat"`
	Plon float64 `bson:"lon"`
}

// Implement Point interface
func (p *Point) Lat() float64 { return p.Plat }
func (p *Point) Lon() float64 { return p.Plon }
func (p *Point) Id() string   { return p.Pid }

// newGeoIndex returns a geo index with points loaded
func newGeoIndex() *geoindex.ClusteringIndex {

	//s := session.Copy()
        //defer s.Close()
        f, _ := os.Open("geo_db.json")
        c, _ := ioutil.ReadAll(f)

	points := make([]*Point, 0)
	err := json.Unmarshal(c, &points)
	if err != nil {
		log.Println("Failed get geo data: ", err)
	}

	fmt.Printf("newGeoIndex len(points) = %d\n", len(points))

	// add points to index
	index := geoindex.NewClusteringIndex()
	for _, point := range points {
		index.Add(point)
	}

	return index
}

func getNearbyPoints(lat, lon float64) []geoindex.Point {
	// fmt.Printf("In geo getNearbyPoints, lat = %f, lon = %f\n", lat, lon)

	center := &geoindex.GeoPoint{
		Pid:  "",
		Plat: lat,
		Plon: lon,
	}

	//MongoSession, _ := mgo.Dial(os.Getenv("HOTEL_APP_DATABASE"))
	return newGeoIndex().KNearest(
		center,
		maxSearchResults,
		geoindex.Km(maxSearchRadius), func(p geoindex.Point) bool {
			return true
		},
	)
}

// Nearby returns all hotels within a given distance.
func Nearby(req RequestBody) string {
	// fmt.Printf("In geo Nearby\n")

	var (
		points = getNearbyPoints(float64(req.Lat), float64(req.Lon))
		res    = make([]string, 0)
	)

	// fmt.Printf("geo after getNearbyPoints, len = %d\n", len(points))

	for _, p := range points {
		// fmt.Printf("In geo Nearby return hotelId = %s\n", p.Id())
		res = append(res, p.Id())
	}

	ret, _ := json.Marshal(res)
	return string(ret)
}

// Handle an HTTP Request.
func FunctionHandler(res http.ResponseWriter, req RequestBody, content_type string, is_json bool) {
	body_u := req
	workflow_id := body_u.WorkflowID
        workflow_depth := body_u.WorkflowDepth
        workflow_width := body_u.WorkflowWidth
        body_u.WorkflowDepth += 1
	fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "3")
	ret := Nearby(body_u)
	fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "4")
	fmt.Fprintf(res, ret) // echo to caller
}
