package function

import (
	"context"
	"fmt"
	"net/http"
	"math"
	"os"
	"encoding/json"
	"io/ioutil"
        "strconv"
        "math/rand"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/sirupsen/logrus"

	"time"

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

type Hotel struct {
	ID     bson.ObjectId `bson:"_id"`
	HId    string        `bson:"hotelid"`
	HLat   float64       `bson:"lat"`
	HLon   float64       `bson:"lon"`
	HRate  float64       `bson:"rate"`
	HPrice float64       `bson:"price"`
}

// loadRecommendations loads hotel recommendations from mongodb.
func loadRecommendations(session *mgo.Session) map[string]Hotel {
	// session, err := mgo.Dial("mongodb-recommendation")
	// if err != nil {
	// 	panic(err)
	// }
	// defer session.Close()
	s := session.Copy()
	defer s.Close()

	c := s.DB("recommendation-db").C("recommendation")

	// unmarshal json profiles
	var hotels []Hotel
	err := c.Find(bson.M{}).All(&hotels)
	if err != nil {
		log.Println("Failed get hotels data: ", err)
	}

	profiles := make(map[string]Hotel)
	for _, hotel := range hotels {
		profiles[hotel.HId] = hotel
	}

	return profiles
}

// GiveRecommendation returns recommendations within a given requirement.
func GetRecommendations(req RequestBody) string {
	MongoSession, _ := mgo.Dial(os.Getenv("HOTEL_APP_DATABASE"))
	var hotels map[string]Hotel
	hotels = loadRecommendations(MongoSession)

	res := make([]string, 0)
	// fmt.Printf("GetRecommendations:\n")
	// fmt.Printf("%+v\n", s.hotels)

	require := req.Require
	if require == "dis" {
		p1 := &geoindex.GeoPoint{
			Pid:  "",
			Plat: req.Lat,
			Plon: req.Lon,
		}
		min := math.MaxFloat64
		for _, hotel := range hotels {
			tmp := float64(geoindex.Distance(p1, &geoindex.GeoPoint{
				Pid:  "",
				Plat: hotel.HLat,
				Plon: hotel.HLon,
			})) / 1000
			if tmp < min {
				min = tmp
			}
		}
		for _, hotel := range hotels {
			tmp := float64(geoindex.Distance(p1, &geoindex.GeoPoint{
				Pid:  "",
				Plat: hotel.HLat,
				Plon: hotel.HLon,
			})) / 1000
			if tmp == min {
				res = append(res, hotel.HId)
			}
		}
	} else if require == "rate" {
		max := 0.0
		for _, hotel := range hotels {
			if hotel.HRate > max {
				max = hotel.HRate
			}
		}
		for _, hotel := range hotels {
			if hotel.HRate == max {
				res = append(res, hotel.HId)
			}
		}
	} else if require == "price" {
		min := math.MaxFloat64
		for _, hotel := range hotels {
			if hotel.HPrice < min {
				min = hotel.HPrice
			}
		}
		for _, hotel := range hotels {
			if hotel.HPrice == min {
				res = append(res, hotel.HId)
			}
		}
	} else {
		log.Println("Wrong parameter: ", require)
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
                body_u.WorkflowDepth += 1
        } else {
                body_u.WorkflowID = workflow_id
                body_u.WorkflowDepth = workflow_depth
                body_u.WorkflowWidth = workflow_width
        }
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "0")
	ret := GetRecommendations(body_u)
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "1")
	fmt.Fprintf(res, ret) // echo to caller
}
