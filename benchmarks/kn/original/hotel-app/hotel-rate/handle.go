package function

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"encoding/json"
	"sort"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"io/ioutil"
        "strconv"
        "math/rand"

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
	WorkflowID string `json:"WorkflowID"`
        WorkflowDepth int `json:"WorkflowDepth"`
        WorkflowWidth int `json:"WorkflowWidth"`
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

func (r RatePlans) Len() int {
	return len(r)
}

func (r RatePlans) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RatePlans) Less(i, j int) bool {
	return r[i].roomType.totalRate > r[j].roomType.totalRate
}

// GetRates gets rates for hotels for specific date range.
func GetRates(req RequestBody) string {
	var res RatePlans
	// session, err := mgo.Dial("mongodb-rate")
	// if err != nil {
	// 	panic(err)
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

			var tmpRatePlans RatePlans

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
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "0" + "\n")
	ret := GetRates(body_u)
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "1" + "\n")
	fmt.Fprintf(res, ret) // echo to caller
}
