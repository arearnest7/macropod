package function

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"bytes"
	"strings"
	"encoding/json"
	"io/ioutil"
        "sort"
        "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/hailocab/go-geoindex"

        log "github.com/sirupsen/logrus"

	"time"

	"github.com/bradfitz/gomemcache/memcache"
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

const (
        maxSearchRadius  = 10
        maxSearchResults = 5
)

type Point struct {
        Pid  string  `bson:"hotelid"`
        Plat float64 `bson:"lat"`
        Plon float64 `bson:"lon"`
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

type User struct {
        Username string `bson2:"username" json:"username"`
        Password string `bson2:"password" json:"password"`
}

type Reservation struct {
        HotelId      string `bson:"hotelid"`
        CustomerName string `bson:"customername"`
        InDate       string `bson:"indate"`
        OutDate      string `bson:"outdate"`
        Number       int    `bson:"number"`
}

type Number struct {
        HotelId string `bson:"hotelid"`
        Number  int    `bson:"numberofroom"`
}

type Hotel struct {
        ID     bson.ObjectId `bson:"_id"`
        HId    string        `bson:"hotelid"`
        HLat   float64       `bson:"lat"`
        HLon   float64       `bson:"lon"`
        HRate  float64       `bson:"rate"`
        HPrice float64       `bson:"price"`
}

type Image struct {
        url string
        Default bool
}

type Address struct {
        streetNumber string
        streetName string
        city string
        state string
        country string
        postalCode string
        lat float64
        lon float64
}

type Hotel2 struct {
        id string
        name string
        phoneNumber string
        description string
        address Address
        images []Image
}

func (r RatePlans) Len() int {
        return len(r)
}

func (r RatePlans) Swap(i, j int) {
        r[i], r[j] = r[j], r[i]
}

func (r RatePlans) Less(i, j int) bool {
        return r[i].roomType.totalRate > r[j].roomType.totalRate
}

// Implement Point interface
func (p *Point) Lat() float64 { return p.Plat }
func (p *Point) Lon() float64 { return p.Plon }
func (p *Point) Id() string   { return p.Pid }

// GetProfiles returns hotel profiles for requested IDs
func GetProfiles(req RequestBody) string {
        // session, err := mgo.Dial("mongodb-profile")
        // if err != nil {
        //      panic(err)
        // }
        // defer session.Close()
        // fmt.Printf("In GetProfiles\n")

        // fmt.Printf("In GetProfiles after setting c\n")

        res := make([]Hotel2, 0)
        hotels := make([]Hotel2, 0)

        MongoSession, _ := mgo.Dial(os.Getenv("HOTEL_APP_DATABASE"))
        var MemcClient = memcache.New(os.Getenv("HOTEL_APP_MEMCACHED"))

	// one hotel should only have one profile

        for _, i := range req.HotelIds {
                // first check memcached
                item, err := MemcClient.Get(i)
                if err == nil {
                        // memcached hit
                        // profile_str := string(item.Value)

                        // fmt.Printf("memc hit\n")
                        // fmt.Println(profile_str)

                        hotel_prof := new(Hotel2)
                        if err = json.Unmarshal(item.Value, hotel_prof); err != nil {
                                log.Warn(err)
                        }
                        hotels = append(hotels, *hotel_prof)

                } else if err == memcache.ErrCacheMiss {
                        // memcached miss, set up mongo connection
                        session := MongoSession.Copy()
                        defer session.Close()
                        c := session.DB("profile-db").C("hotels")

                        hotel_prof := new(Hotel2)
                        err := c.Find(bson.M{"id": i}).One(&hotel_prof)

                        if err != nil {
                                log.Println("Hotel data not found in memcached: ", err)
                        }

                        // for _, h := range hotels {
			//      res.Hotels = append(res.Hotels, h)
                        // }
                        hotels = append(hotels, *hotel_prof)

                        prof_json, err := json.Marshal(hotel_prof)
                        if err != nil {
                                log.Warn(err)
                        }
                        memc_str := string(prof_json)

                        // write to memcached
                        err = MemcClient.Set(&memcache.Item{Key: i, Value: []byte(memc_str)})
                        if err != nil {
                                log.Warn("MMC error: ", err)
                        }
                } else {
                        fmt.Printf("Memmcached error = %s\n", err)
                        panic(err)
                }
        }

        res = hotels
        // fmt.Printf("In GetProfiles after getting resp\n")
        ret, _ := json.Marshal(res)
        return string(ret)
}

// GetRates gets rates for hotels for specific date range.
func GetRates(req Request) string {
        var res RatePlans
        // session, err := mgo.Dial("mongodb-rate")
        // if err != nil {
        //      panic(err)
        // }
        // defer session.Close()

        ratePlans := make(RatePlans, 0)

        MongoSession, _ := mgo.Dial(os.Getenv("HOTEL_APP_DATABASE"))
        var MemcClient = memcache.New(os.Getenv("HOTEL_APP_MEMCACHED"))

        // fmt.Printf("Hotel Ids: %+v\n", req.HotelIds)

        for _, hotelID := range req.HotelIds {
                // first check memcached
                item, err := MemcClient.Get(hotelID)
                if err == nil {
                        // memcached hit
                        rate_strs := strings.Split(string(item.Value), "\n")

                        // fmt.Printf("memc hit, hotelId = %s\n", hotelID)
                        fmt.Println(rate_strs)

                        for _, rate_str := range rate_strs {
                                if len(rate_str) != 0 {
                                        rate_p := new(RatePlan)
                                        if err = json.Unmarshal(item.Value, rate_p); err != nil {
                                                log.Warn(err)
                                        }
                                        ratePlans = append(ratePlans, rate_p)
                                }
                        }
                } else if err == memcache.ErrCacheMiss {

                        // fmt.Printf("memc miss, hotelId = %s\n", hotelID)

                        // memcached miss, set up mongo connection
                        session := MongoSession.Copy()
                        defer session.Close()
                        c := session.DB("rate-db").C("inventory")

                        memc_str := ""

                        tmpRatePlans := make(RatePlans, 0)

                        err = c.Find(&bson.M{"hotelid": hotelID}).All(&tmpRatePlans)
                        // fmt.Printf("Rate Plans %+v\n", tmpRatePlans)
			if err != nil {
                                panic(err)
                        } else {
                                for _, r := range tmpRatePlans {
                                        ratePlans = append(ratePlans, r)
                                        rate_json, err := json.Marshal(r)
                                        if err != nil {
                                                fmt.Printf("json.Marshal err = %s\n", err)
                                        }
                                        memc_str = memc_str + string(rate_json) + "\n"
                                }
                        }

                        // write to memcached
                        err = MemcClient.Set(&memcache.Item{Key: hotelID, Value: []byte(memc_str)})
                        if err != nil {
                                log.Warn("MMC error: ", err)
                        }
                } else {
                        fmt.Printf("Memmcached error = %s\n", err)
                        panic(err)
                }
        }

        // fmt.Printf("Rate Plans %+v\n", ratePlans)
        sort.Sort(ratePlans)
        res = ratePlans

        ret, _ := json.Marshal(res)
        return string(ret)
}

// newGeoIndex returns a geo index with points loaded
func newGeoIndex(session *mgo.Session) *geoindex.ClusteringIndex {

        s := session.Copy()
        defer s.Close()
        c := s.DB("geo-db").C("geo")

        points := make([]*Point, 0)
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
func Nearby(req BodyGeo) string {
        // fmt.Printf("In geo Nearby\n")

        var (
                points = getNearbyPoints(req.Lat, req.Lon)
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

// Nearby returns ids of nearby hotels ordered by ranking algo
func SearchNearby(req RequestBody) string {
        // find nearby hotels
        fmt.Printf("in Search Nearby\n")

        fmt.Printf("nearby lat = %f\n", req.Lat)
        fmt.Printf("nearby lon = %f\n", req.Lon)

        payload := BodyGeo{Lat: req.Lat, Lon: req.Lon}
	nearby := Nearby(payload)
        //if err != nil {
        //        fmt.Printf("nearby error: %v", err)
        //        return ""
        //}

        // var ids []string
        nearby_u := make([]string, 0)
        _ = json.Unmarshal([]byte(nearby), &nearby_u)
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

        rates := GetRates(r)
        //if err != nil {
        //        fmt.Printf("rates error: %v", err)
        //        return ""
        //}
        // TODO(hw): add simple ranking algo to order hotel ids:
	// * geo distance
        // * price (best discount?)
        // * reviews

        // build the response
        res := make([]string, 0)
        var rate_p RatePlans
        json.Unmarshal([]byte(rates), &rate_p)
        for _, ratePlan := range rate_p {
                // fmt.Printf("get RatePlan HotelId = %s, Code = %s\n", ratePlan.HotelId, ratePlan.Code)
                res = append(res, ratePlan.hotelId)
        }
        ret, _ := json.Marshal(res)
        return string(ret)
}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
        fmt.Println(time.Now().String() + "," + "0" + "," + "0" + "," + "0" + "," + "HTTP" + "," + "0" + "\n")
	ret := ""
        body, _ := ioutil.ReadAll(req.Body)
        body_u := RequestBody{}
        json.Unmarshal(body, &body_u)
        defer req.Body.Close()
        if body_u.Request == "search" {
                ret = SearchNearby(body_u)
        } else {
		requestURL := ""
                if body_u.Request == "recommend" {
                        requestURL = os.Getenv("HOTEL_RECOMMEND_PARTIAL")
                } else if body_u.Request == "reserve" {
                        requestURL = os.Getenv("HOTEL_RESERVE_PARTIAL")
                } else if body_u.Request == "user" {
                        requestURL = os.Getenv("HOTEL_USER_PARTIAL")
                }

                body_m, err := json.Marshal(body_u)
        	req_url, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(body_m))
        	if err != nil {
                	log.Fatal(err)
        	}
        	req_url.Header.Add("Content-Type", "application/json")
        	client := &http.Client{}
                fmt.Println(time.Now().String() + "," + "0" + "," + "0" + "," + "0" + "," + "HTTP" + "," + "1" + "\n")
        	ret1, err := client.Do(req_url)
                fmt.Println(time.Now().String() + "," + "0" + "," + "0" + "," + "0" + "," + "HTTP" + "," + "2" + "\n")
        	retBody, err := ioutil.ReadAll(ret1.Body)
        	ret_val, err := json.Marshal(retBody)
		ret = string(ret_val)
        }
        fmt.Println(time.Now().String() + "," + "0" + "," + "0" + "," + "0" + "," + "HTTP" + "," + "3" + "\n")
        fmt.Fprintf(res, ret) // echo to caller
}

