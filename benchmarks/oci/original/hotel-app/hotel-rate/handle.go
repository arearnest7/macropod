package function

import (
	"fmt"
	"net/http"
	"encoding/json"
	"sort"
	"os"
	"io/ioutil"
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

	//MongoSession, _ := mgo.Dial(os.Getenv("HOTEL_APP_DATABASE"))
        //var MemcClient = memcache.New(os.Getenv("HOTEL_APP_MEMCACHED"))

	// fmt.Printf("Hotel Ids: %+v\n", req.HotelIds)

	for _, hotelID := range req.HotelIds {
		// first check memcached
		_, err := "not nil", memcache.ErrCacheMiss
		if err == nil {
			// memcached hit
			//rate_strs := strings.Split(string(item.Value), "\n")
                        rate_strs := make([]string, 0)
			// fmt.Printf("memc hit, hotelId = %s\n", hotelID)
			//fmt.Println(rate_strs)

			for _, rate_str := range rate_strs {
				if len(rate_str) != 0 {
					rate_p := new(RatePlan)
					//if err = json.Unmarshal(item.Value, rate_p); err != nil {
					//	log.Warn(err)
					//}
					ratePlans = append(ratePlans, rate_p)
				}
			}
		} else if err == memcache.ErrCacheMiss {

			// fmt.Printf("memc miss, hotelId = %s\n", hotelID)

			// memcached miss, set up mongo connection
			//session := MongoSession.Copy()
                        //defer session.Close()
                        f, _ := os.Open("rate_db.json")
                        c, _ := ioutil.ReadAll(f)

			memc_str := ""

			tmpRatePlans := make(RatePlans, 0)
                        tmpRatePlans_temp := make(RatePlans, 0)
                        err := json.Unmarshal(c, &tmpRatePlans_temp)
                        for _, h := range tmpRatePlans_temp {
                                if h.hotelId == hotelID {
                                        tmpRatePlans = append(tmpRatePlans, h)
                                }
                        }
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
			//err = memcache.ErrCacheMiss
			//if err != nil {
			//	log.Warn("MMC error: ", err)
			//}
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
func FunctionHandler(res http.ResponseWriter, req *http.Request, content_type string, is_json bool) {
	fmt.Println(time.Now().String() + "," + "0" + "," + "0" + "," + "0" + "," + "HTTP" + "," + "2" + "\n")
        body, _ := ioutil.ReadAll(req.Body)
        body_u := RequestBody{}
        json.Unmarshal(body, &body_u)
        defer req.Body.Close()
	fmt.Println(time.Now().String() + "," + "0" + "," + "0" + "," + "0" + "," + "HTTP" + "," + "3" + "\n")
	fmt.Fprintf(res, GetRates(body_u)) // echo to caller
}
