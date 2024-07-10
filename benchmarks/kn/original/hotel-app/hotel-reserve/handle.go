package function

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"os"
	"encoding/json"
	"io/ioutil"
        "math/rand"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

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

// CheckAvailability checks if given information is available
func CheckAvailability(req RequestBody) string {
	log.Println("CheckAvailability")
	res := make([]string, 0)

	// session, err := mgo.Dial("mongodb-reservation")
	// if err != nil {
	// 	panic(err)
	// }
	// defer session.Close()
	MongoSession, _ := mgo.Dial(os.Getenv("HOTEL_APP_DATABASE"))
        var MemcClient = memcache.New(os.Getenv("HOTEL_APP_MEMCACHED"))
	session := MongoSession.Copy()
	defer session.Close()

	c := session.DB("reservation-db").C("reservation")
	c1 := session.DB("reservation-db").C("number")

	for _, hotelId := range req.HotelIds {
		fmt.Printf("reservation check hotel %s\n", hotelId)
		inDate, _ := time.Parse(
			time.RFC3339,
			req.InDate+"T12:00:00+00:00")

		outDate, _ := time.Parse(
			time.RFC3339,
			req.OutDate+"T12:00:00+00:00")

		indate := inDate.String()[0:10]

		for inDate.Before(outDate) {
			// check reservations
			count := 0
			inDate = inDate.AddDate(0, 0, 1)
			fmt.Printf("reservation check date %s\n", inDate.String()[0:10])
			outdate := inDate.String()[0:10]

			// first check memc
			memc_key := hotelId + "_" + inDate.String()[0:10] + "_" + outdate
			item, err := MemcClient.Get(memc_key)

			if err == nil {
				// memcached hit
				count, _ = strconv.Atoi(string(item.Value))
				fmt.Printf("memcached hit %s = %d\n", memc_key, count)
			} else if err == memcache.ErrCacheMiss {
				// memcached miss
				reserve := make([]Reservation, 0)
				err := c.Find(&bson.M{"hotelid": hotelId, "inDate": indate, "outDate": outdate}).All(&reserve)
				if err != nil {
					panic(err)
				}
				for _, r := range reserve {
					fmt.Printf("reservation check reservation number = %s\n", hotelId)
					count += r.Number
				}

				// update memcached
				err = MemcClient.Set(&memcache.Item{Key: memc_key, Value: []byte(strconv.Itoa(count))})
				if err != nil {
					log.Warn("MMC error: ", err)
				}
			} else {
				fmt.Printf("Memmcached error = %s\n", err)
				panic(err)
			}

			// check capacity
			// check memc capacity
			memc_cap_key := hotelId + "_cap"
			item, err = MemcClient.Get(memc_cap_key)
			hotel_cap := 0

			if err == nil {
				// memcached hit
				hotel_cap, _ = strconv.Atoi(string(item.Value))
				fmt.Printf("memcached hit %s = %d\n", memc_cap_key, hotel_cap)
			} else if err == memcache.ErrCacheMiss {
				var num Number
				err = c1.Find(&bson.M{"hotelid": hotelId}).One(&num)
				if err != nil {
					panic(err)
				}
				hotel_cap = int(num.Number)
				// update memcached
				err = MemcClient.Set(&memcache.Item{Key: memc_cap_key, Value: []byte(strconv.Itoa(hotel_cap))})
				if err != nil {
					log.Warn("MMC error: ", err)
				}
			} else {
				fmt.Printf("Memmcached error = %s\n", err)
				panic(err)
			}

			if count+int(req.RoomNumber) > hotel_cap {
				break
			}
			indate = outdate

			if inDate.Equal(outDate) {
				res = append(res, hotelId)
			}
		}
	}

	ret, _ := json.Marshal(res)
        return string(ret)
}

// MakeReservation makes a reservation based on given information
func MakeReservation(req RequestBody) string {
	log.Println("MakeReservation")
	res := make([]string, 0)

	// session, err := mgo.Dial("mongodb-reservation")
	// if err != nil {
	// 	panic(err)
	// }
	// defer session.Close()
	MongoSession, _ := mgo.Dial(os.Getenv("HOTEL_APP_DATABASE"))
        var MemcClient = memcache.New(os.Getenv("HOTEL_APP_MEMCACHED"))
	session := MongoSession.Copy()
	defer session.Close()

	c := session.DB("reservation-db").C("reservation")
	c1 := session.DB("reservation-db").C("number")

	inDate, _ := time.Parse(
		time.RFC3339,
		req.InDate+"T12:00:00+00:00")

	outDate, _ := time.Parse(
		time.RFC3339,
		req.OutDate+"T12:00:00+00:00")
	hotelId := req.HotelId

	indate := inDate.String()[0:10]

	memc_date_num_map := make(map[string]int)

	for inDate.Before(outDate) {
		// check reservations
		count := 0
		inDate = inDate.AddDate(0, 0, 1)
		outdate := inDate.String()[0:10]

		// first check memc
		memc_key := hotelId + "_" + inDate.String()[0:10] + "_" + outdate
		item, err := MemcClient.Get(memc_key)
		if err == nil {
			// memcached hit
			count, _ = strconv.Atoi(string(item.Value))
			fmt.Printf("memcached hit %s = %d\n", memc_key, count)
			memc_date_num_map[memc_key] = count + int(req.RoomNumber)

		} else if err == memcache.ErrCacheMiss {
			// memcached miss
			fmt.Printf("memcached miss\n")
			reserve := make([]Reservation, 0)
			err := c.Find(&bson.M{"hotelid": hotelId, "inDate": indate, "outDate": outdate}).All(&reserve)
			if err != nil {
				panic(err)
			}

			for _, r := range reserve {
				count += r.Number
			}

			memc_date_num_map[memc_key] = count + int(req.RoomNumber)

		} else {
			fmt.Printf("Memmcached error = %s\n", err)
			panic(err)
		}

		// check capacity
		// check memc capacity
		memc_cap_key := hotelId + "_cap"
		item, err = MemcClient.Get(memc_cap_key)
		hotel_cap := 0
		if err == nil {
			// memcached hit
			hotel_cap, _ = strconv.Atoi(string(item.Value))
			fmt.Printf("memcached hit %s = %d\n", memc_cap_key, hotel_cap)
		} else if err == memcache.ErrCacheMiss {
			// memcached miss
			var num Number
			err = c1.Find(&bson.M{"hotelid": hotelId}).One(&num)
			if err != nil {
				panic(err)
			}
			hotel_cap = int(num.Number)

			// write to memcache
			err = MemcClient.Set(&memcache.Item{Key: memc_cap_key, Value: []byte(strconv.Itoa(hotel_cap))})
			if err != nil {
				log.Warn("MMC error: ", err)
			}
		} else {
			fmt.Printf("Memmcached error = %s\n", err)
			panic(err)
		}

		if count+int(req.RoomNumber) > hotel_cap {
			fmt.Printf("Not enough space left\n")
			ret, _ := json.Marshal(res)
        		return string(ret)
		}
		indate = outdate
	}

	// only update reservation number cache after check succeeds
	for key, val := range memc_date_num_map {
		err := MemcClient.Set(&memcache.Item{Key: key, Value: []byte(strconv.Itoa(val))})
		if err != nil {
			log.Warn("MMC error: ", err)
		}
	}

	inDate, _ = time.Parse(
		time.RFC3339,
		req.InDate+"T12:00:00+00:00")

	indate = inDate.String()[0:10]

	for inDate.Before(outDate) {
		inDate = inDate.AddDate(0, 0, 1)
		outdate := inDate.String()[0:10]
		err := c.Insert(&Reservation{
			HotelId:      hotelId,
			CustomerName: req.CustomerName,
			InDate:       indate,
			OutDate:      outdate,
			Number:       int(req.RoomNumber)})
		if err != nil {
			panic(err)
		}
		indate = outdate
	}

	res = append(res, hotelId)

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
        ret := ""
	if body_u.RequestType == "check" {
		ret = CheckAvailability(body_u)
	} else if body_u.RequestType == "make" {
		ret = MakeReservation(body_u)
	}
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "1")
	fmt.Fprintf(res, ret) // echo to caller
}
