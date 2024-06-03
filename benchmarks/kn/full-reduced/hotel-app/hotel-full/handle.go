package function

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"encoding/json"
	"io/ioutil"
        "sort"
        "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math"
	"strconv"
        "math/rand"
	"crypto/sha256"

	bson2 "go.mongodb.org/mongo-driver/bson"
        "go.mongodb.org/mongo-driver/bson/primitive"
        "go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"

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
	WorkflowID string `json:"WorkflowID"`
        WorkflowDepth int `json:"WorkflowDepth"`
        WorkflowWidth int `json:"WorkflowWidth"`
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

// loadRecommendations loads hotel recommendations from mongodb.
func loadRecommendations(session *mgo.Session) map[string]Hotel {
        // session, err := mgo.Dial("mongodb-recommendation")
        // if err != nil {
        //      panic(err)
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

// CheckAvailability checks if given information is available
func CheckAvailability(req RequestBody) string {
        log.Println("CheckAvailability")
        res := make([]string, 0)

        // session, err := mgo.Dial("mongodb-reservation")
        // if err != nil {
        //      panic(err)
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
                                err = MemcClient.Set(&memcache.Item{Key: memc_cap_key, Value: []byte(strconv.Itoa(count))})
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
        //      panic(err)
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

// loadUsers loads hotel users from database
func loadUsers(client *mongo.Client) map[string]string {

        coll := client.Database("user-db").Collection("user")

        filter := bson2.D{}
        // filter := bson.M{{"username": username}}
        cursor, err := coll.Find(context.Background(), filter)
        if err != nil {
                log.Println("Failed get users data: ", err)
        }

        // Get a list of all returned documents and print them out.
        // See the mongo.Cursor documentation for more examples of using cursors.
        var users []User
        if err = cursor.All(context.Background(), &users); err != nil {
                log.Println("Failed get users data: ", err)
        }

        res := make(map[string]string)
        for _, user := range users {
                res[user.Username] = user.Password
        }

        fmt.Printf("Done load users\n")

        return res
}

func lookupCache(username string) string {
        ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
        defer cancel()
        MongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("HOTEL_APP_DATABASE")))
        if err != nil { return "" }
        users_cached := loadUsers(MongoClient)
        res, ok := users_cached[username]
        if !ok {
                log.Println("User does not exist: ", username)
        }
        return res
}

//
func lookUpDB(username string) (User, bool) {
        // session := s.MongoClient.Copy()
        // defer session.Close()
        ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
        defer cancel()
        MongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("HOTEL_APP_DATABASE")))
        collection := MongoClient.Database("user-db").Collection("user")

        // listAll(collection)

        // unmarshal json profiles
        var user User
        filter := bson2.D{primitive.E{Key: "username", Value: username}}
        // filter := bson.M{{"username": username}}
        err = collection.FindOne(context.Background(), filter).Decode(&user)
        if err != nil {
                log.Println("Failed get user: ", err)
                return user, false
        }
        return user, true
}

// CheckUser returns whether the username and password are correct.
func CheckUser(req RequestBody) bool {
        var res bool

        fmt.Printf("CheckUser: %+v", req)

        sum := sha256.Sum256([]byte(req.Password))
        pass := fmt.Sprintf("%x", sum)

        use_cache := false

        if use_cache {
                password := lookupCache(req.Username)
                res = pass == password
        } else {
                user, _ := lookUpDB(req.Username)
                res = pass == user.Password
        }

        fmt.Printf(" >> pass: %t\n", res)

        return res
}

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
	ret := ""
	if body_u.Request == "search" {
		ret = SearchNearby(body_u)
	} else if body_u.Request == "recommend" {
		ret = GetRecommendations(body_u)
	} else if body_u.Request == "reserve" {
		if body_u.RequestType == "check" {
			ret = CheckAvailability(body_u)
		} else if body_u.RequestType == "make" {
			ret = MakeReservation(body_u)
		}
	} else if body_u.Request == "user" {
		ret = strconv.FormatBool(CheckUser(body_u))
	}
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + "," + workflow_id + "," + strconv.Itoa(workflow_depth) + "," + strconv.Itoa(workflow_width) + "," + "HTTP" + "," + "1" + "\n")
	fmt.Fprintf(res, ret) // echo to caller
}

