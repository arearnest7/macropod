package function

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"encoding/json"
	"io/ioutil"

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
	//MongoSession, _ := mgo.Dial(os.Getenv("HOTEL_APP_DATABASE"))
        //var MemcClient = memcache.New(os.Getenv("HOTEL_APP_MEMCACHED"))
        //session := MongoSession.Copy()
        //defer session.Close()

        f, _ := os.Open("reservation_db.json")
        c1, _ := ioutil.ReadAll(f)
        c := []byte("{}")

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
			_, err := "not nil", memcache.ErrCacheMiss

			if err == nil {
				// memcached hit
				//count, _ = strconv.Atoi(string(item.Value))
				fmt.Printf("memcached hit %s = %d\n", memc_key, "")
			} else if err == memcache.ErrCacheMiss {
				// memcached miss
				reserve := make([]Reservation, 0)
				temp_reserve := make([]Reservation, 0)
                                err := json.Unmarshal(c1, &temp_reserve)
                                for _, r := range temp_reserve {
                                        if r.HotelId == hotelId && r.InDate == indate && r.OutDate == outdate {
                                                reserve = append(reserve, r)
                                        }
                                }
				if err != nil {
					panic(err)
				}
				for _, r := range reserve {
					fmt.Printf("reservation check reservation number = %s\n", hotelId)
					count += r.Number
				}

				// update memcached
				//err = MemcClient.Set(&memcache.Item{Key: memc_key, Value: []byte(strconv.Itoa(count))})
				//if err != nil {
				//	log.Warn("MMC error: ", err)
				//}
			} else {
				fmt.Printf("Memmcached error = %s\n", err)
				panic(err)
			}

			// check capacity
			// check memc capacity
			memc_cap_key := hotelId + "_cap"
			_, err = "not nil", memcache.ErrCacheMiss
			hotel_cap := 0

			if err == nil {
				// memcached hit
				//hotel_cap, _ = strconv.Atoi(string(item.Value))
				fmt.Printf("memcached hit %s = %d\n", memc_cap_key, "")
			} else if err == memcache.ErrCacheMiss {
				var num Number
				nums := make([]Number, 0)
                                err := json.Unmarshal(c, &nums)
                                for _, n := range nums {
                                        if n.HotelId == hotelId {
                                                num = n
                                        }
                                }
				if err != nil {
					panic(err)
				}
				hotel_cap = int(num.Number)
				// update memcached
				//err = MemcClient.Set(&memcache.Item{Key: memc_cap_key, Value: []byte(strconv.Itoa(hotel_cap))})
				//if err != nil {
				//	log.Warn("MMC error: ", err)
				//}
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
	//MongoSession, _ := mgo.Dial(os.Getenv("HOTEL_APP_DATABASE"))
        //var MemcClient = memcache.New(os.Getenv("HOTEL_APP_MEMCACHED"))
        //session := MongoSession.Copy()
        //defer session.Close()

        f, _ := os.Open("reservation_db.json")
        c1, _ := ioutil.ReadAll(f)
        c := []byte("{}")

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
		_, err := "not nil", memcache.ErrCacheMiss
		if err == nil {
			// memcached hit
			//count, _ = strconv.Atoi(string(item.Value))
			fmt.Printf("memcached hit %s = %d\n", memc_key, "")
			//memc_date_num_map[memc_key] = count + int(req.RoomNumber)

		} else if err == memcache.ErrCacheMiss {
			// memcached miss
			fmt.Printf("memcached miss\n")
			reserve := make([]Reservation, 0)
			temp_reserve := make([]Reservation, 0)
                        err := json.Unmarshal(c1, &temp_reserve)
                        for _, r := range temp_reserve {
                                if r.HotelId == hotelId && r.InDate == indate && r.OutDate == outdate {
                                        reserve = append(reserve, r)
                                }
                        }
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
		_, err = "not nil", memcache.ErrCacheMiss
		hotel_cap := 0
		if err == nil {
			// memcached hit
			//hotel_cap, _ = strconv.Atoi(string(item.Value))
			fmt.Printf("memcached hit %s = %d\n", memc_cap_key, "")
		} else if err == memcache.ErrCacheMiss {
			// memcached miss
			var num Number
			nums := make([]Number, 0)
                        err := json.Unmarshal(c, &nums)
                        for _, n := range nums {
                                if n.HotelId == hotelId {
                                        num = n
                                }
                        }
			if err != nil {
				panic(err)
			}
			hotel_cap = int(num.Number)

			// write to memcache
			//err = MemcClient.Set(&memcache.Item{Key: memc_cap_key, Value: []byte(strconv.Itoa(hotel_cap))})
			//if err != nil {
			//	log.Warn("MMC error: ", err)
			//}
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
	//for key, val := range memc_date_num_map {
	//	err := memcache.ErrCacheMiss
	//	if err != nil {
	//		log.Warn("MMC error: ", err)
	//	}
	//}

	inDate, _ = time.Parse(
		time.RFC3339,
		req.InDate+"T12:00:00+00:00")

	indate = inDate.String()[0:10]

	for inDate.Before(outDate) {
		inDate = inDate.AddDate(0, 0, 1)
		outdate := inDate.String()[0:10]
		//err := c.Insert(&Reservation{
		//	HotelId:      hotelId,
		//	CustomerName: req.CustomerName,
		//	InDate:       indate,
		//	OutDate:      outdate,
		//	Number:       int(req.RoomNumber)})
		//if err != nil {
		//	panic(err)
		//}
		indate = outdate
	}

	res = append(res, hotelId)

	ret, _ := json.Marshal(res)
        return string(ret)
}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	logging_name, logging := os.LookupEnv("LOGGING_NAME")
        redisClient := redis.NewClient(&redis.Options{})
        c := context.Background()
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
        ret := ""
	body, _ := ioutil.ReadAll(req.Body)
        body_u := RequestBody{}
        json.Unmarshal(body, &body_u)
        defer req.Body.Close()
	if body_u.RequestType == "check" {
		ret = CheckAvailability(body_u)
	} else if body_u.RequestType == "make" {
		ret = MakeReservation(body_u)
	}
	if logging {
                redisClient.Append(c, logging_name, time.Now().String() + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "1" + "\n")
        }
	fmt.Fprintf(res, ret) // echo to caller
}
