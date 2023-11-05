package function

import (
	"context"
	"fmt"
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"io/ioutil"
	"encoding/json"


	log "github.com/sirupsen/logrus"

	"github.com/hailocab/go-geoindex"
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
func newGeoIndex(session *mgo.Session) *geoindex.ClusteringIndex {

	s := session.Copy()
	defer s.Close()
	c := s.DB("geo-db").C("geo")

	var points []*Point
	err := c.Find(bson.M{}).All(&points)
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

	MongoSession, _ := mgo.Dial(os.Getenv("HOTEL_APP_DATABASE"))
	return newGeoIndex(MongoSession).KNearest(
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
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
        var body_u *RequestBody
        json.Unmarshal(body, &body_u)
        defer req.Body.Close()
	fmt.Fprintf(res, Nearby(*body_u)) // echo to caller
}
