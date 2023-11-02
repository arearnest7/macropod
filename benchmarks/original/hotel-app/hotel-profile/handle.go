package function

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"encoding/json"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"net"

	log "github.com/sirupsen/logrus"

	"time"

	"github.com/bradfitz/gomemcache/memcache"

	"golang.org/x/net/context"
)

type Image struct {
	url string
	default bool
}

type Address struct {
	streetNumber string
	streetName string
	city string
	state string
	country string
	postalCode string
	lat float
	lon float
}

type Hotel struct {
	id string
	name string
	phoneNumber string
	description string
	address Address
	images Image[]
}

// GetProfiles returns hotel profiles for requested IDs
func GetProfiles(var req) (string, error) {
	// session, err := mgo.Dial("mongodb-profile")
	// if err != nil {
	// 	panic(err)
	// }
	// defer session.Close()
	// fmt.Printf("In GetProfiles\n")

	// fmt.Printf("In GetProfiles after setting c\n")

	res = []Hotel
	hotels := make([]Hotel, 0)

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
			hotels = append(hotels, hotel_prof)

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
			hotels = append(hotels, hotel_prof)

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
	return json.Marhsal(res), nil
}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
        fmt.Fprintf(res, GetProfiles(req.json)) // echo to caller
}
