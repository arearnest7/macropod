package function

import (
    "fmt"
    "encoding/json"
    "sort"
    "os"
    "io/ioutil"

    "github.com/bradfitz/gomemcache/memcache"
)

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
func GetRates(req map[string]interface{}) string {
    var res RatePlans
    // session, err := mgo.Dial("mongodb-rate")
    // if err != nil {
    //     panic(err)
    // }
    // defer session.Close()

    ratePlans := make(RatePlans, 0)

    //MongoSession, _ := mgo.Dial(os.Getenv("HOTEL_APP_DATABASE"))
        //var MemcClient = memcache.New(os.Getenv("HOTEL_APP_MEMCACHED"))

    // fmt.Printf("Hotel Ids: %+v\n", req.HotelIds)

    for _, hotelID := range req["HotelIds"].([]string) {
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
                    //    log.Warn(err)
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
            //    log.Warn("MMC error: ", err)
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

func FunctionHandler(context Context) (string, int) {
    return GetRates(context.JSON), 200
}
