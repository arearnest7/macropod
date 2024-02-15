package function

import (
	"fmt"
	"encoding/json"
	"os"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

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

type Hotel struct {
	id string
	name string
	phoneNumber string
	description string
	address Address
	images []Image
}

// GetProfiles returns hotel profiles for requested IDs
func GetProfiles(req RequestBody) string {
	// session, err := mgo.Dial("mongodb-profile")
	// if err != nil {
	// 	panic(err)
	// }
	// defer session.Close()
	// fmt.Printf("In GetProfiles\n")

	// fmt.Printf("In GetProfiles after setting c\n")

	res := make([]Hotel, 0)
	hotels := make([]Hotel, 0)

	//MongoSession, _ := mgo.Dial(os.Getenv("HOTEL_APP_DATABASE"))
	//var MemcClient = memcache.New(os.Getenv("HOTEL_APP_MEMCACHED"))

	// one hotel should only have one profile

	for _, i := range req.HotelIds {
		// first check memcached
		_, err := "not nil", memcache.ErrCacheMiss
		if err == nil {
			// memcached hit
			// profile_str := string(item.Value)

			// fmt.Printf("memc hit\n")
			// fmt.Println(profile_str)

			hotel_prof := new(Hotel)
			//if err = json.Unmarshal(item.Value, hotel_prof); err != nil {
			//	log.Warn(err)
			//}
			hotels = append(hotels, *hotel_prof)

		} else if err == memcache.ErrCacheMiss {
			// memcached miss, set up mongo connection
			//session := MongoSession.Copy()
                        //defer session.Close()
                        f, _ := os.Open("profile_db.json")
                        c, _ := ioutil.ReadAll(f)

			hotel_prof := new(Hotel)
			hotel_prof_temp := make([]Hotel, 0)
                        err := json.Unmarshal(c, &hotel_prof_temp)
                        for _, h := range hotel_prof_temp {
                                if h.id == i {
                                        hotel_prof = &h
                                }
                        }

			if err != nil {
				log.Println("Hotel data not found in memcached: ", err)
			}

			// for _, h := range hotels {
			// 	res.Hotels = append(res.Hotels, h)
			// }
			hotels = append(hotels, *hotel_prof)

			_, err = json.Marshal(hotel_prof)
			if err != nil {
				log.Warn(err)
			}
			//memc_str := string(prof_json)

			// write to memcached
			//err = MemcClient.Set(&memcache.Item{Key: i, Value: []byte(memc_str)})
			//if err != nil {
			//	log.Warn("MMC error: ", err)
			//}
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

func function_handler(context Context) (string, int) {
        //body, _ := ioutil.ReadAll(req.Body)
        body_u := RequestBody{}
        json.Unmarshal([]byte(context.request), &body_u)
        //defer req.Body.Close()
	return GetProfiles(body_u), 200
}
