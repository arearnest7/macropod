package function

import (
	"context"
	"fmt"
	"net/http"
	"encoding/json"
	"os"
	"io/ioutil"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/sirupsen/logrus"

	"time"
        "github.com/redis/go-redis/v9"

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

			hotel_prof := new(Hotel)
			if err = json.Unmarshal(item.Value, hotel_prof); err != nil {
				log.Warn(err)
			}
			hotels = append(hotels, *hotel_prof)

		} else if err == memcache.ErrCacheMiss {
			// memcached miss, set up mongo connection
			session := MongoSession.Copy()
			defer session.Close()
			c := session.DB("profile-db").C("hotels")

			hotel_prof := new(Hotel)
			err := c.Find(bson.M{"id": i}).One(&hotel_prof)

			if err != nil {
				log.Println("Hotel data not found in memcached: ", err)
			}

			// for _, h := range hotels {
			// 	res.Hotels = append(res.Hotels, h)
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

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	logging_name, logging := os.LookupEnv("LOGGING_NAME")
        redisClient := redis.NewClient(&redis.Options{})
        c := context.Background()
        body, _ := ioutil.ReadAll(req.Body)
        if logging {
                logging_url := os.Getenv("LOGGING_URL")
                logging_password := os.Getenv("LOGGING_PASSWORD")
                redisClient = redis.NewClient(&redis.Options{
                        Addr: logging_url,
                        Password: logging_password,
                        DB: 0,
                })
        }
        if logging {
                redisClient.Append(c, logging_name, time.Now().String() + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n")
        }
        body, _ := ioutil.ReadAll(req.Body)
        body_u := RequestBody{}
        json.Unmarshal(body, &body_u)
        defer req.Body.Close()
	ret := GetProfiles(body_u)
	if logging {
                redisClient.Append(c, logging_name, time.Now().String() + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n")
        }
	fmt.Fprintf(res, ret) // echo to caller
}
