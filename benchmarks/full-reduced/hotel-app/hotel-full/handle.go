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
	"math"
	"strconv"

	"net"

        log "github.com/sirupsen/logrus"

        "time"

	"github.com/bradfitz/gomemcache/memcache"
        "golang.org/x/net/context"
        "google.golang.org/grpc"
        "google.golang.org/grpc/credentials/insecure"
        "google.golang.org/grpc/keepalive"
        "google.golang.org/grpc/reflection"

        geo "github.com/vhive-serverless/vSwarm-proto/proto/hotel_reserv/geo"
        rate "github.com/vhive-serverless/vSwarm-proto/proto/hotel_reserv/rate"

        pb "github.com/vhive-serverless/vSwarm-proto/proto/hotel_reserv/search"

        tracing "github.com/vhive-serverless/vSwarm/utils/tracing/go"
)

const (
        maxSearchRadius  = 10
        maxSearchResults = 5
)

type Point struct {
        Pid  string  `bson:"hotelid"`
        Plat float64 `bson:"lat"`
        Plon float64 `bson:"lon"`
}

type RatePlans []*pb.RatePlan

type User struct {
        Username string `bson:"username" json:"username"`
        Password string `bson:"password" json:"password"`
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

func (r RatePlans) Len() int {
        return len(r)
}

func (r RatePlans) Swap(i, j int) {
        r[i], r[j] = r[j], r[i]
}

func (r RatePlans) Less(i, j int) bool {
        return r[i].RoomType.TotalRate > r[j].RoomType.TotalRate
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
func (s *Server) GetRecommendations(ctx context.Context, req *pb.Request) (*pb.Result, error) {
        map[string]Hotel hotels = loadRecmmendations(MongoSession)

        res := new(pb.Result)
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
                                res.HotelIds = append(res.HotelIds, hotel.HId)
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
                                res.HotelIds = append(res.HotelIds, hotel.HId)
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
                                res.HotelIds = append(res.HotelIds, hotel.HId)
                        }
                }
        } else {
                log.Println("Wrong parameter: ", require)
        }

        return res, nil
}

// CheckAvailability checks if given information is available
func CheckAvailability(var req) (*pb.Result, error) {
        log.Println("CheckAvailability")
        res := new(pb.Result)
        res.HotelId = make([]string, 0)

        // session, err := mgo.Dial("mongodb-reservation")
        // if err != nil {
        //      panic(err)
        // }
        // defer session.Close()
        session := MongoSession.Copy()
        defer session.Close()

        c := session.DB("reservation-db").C("reservation")
        c1 := session.DB("reservation-db").C("number")

        for _, hotelId := range req.HotelId {
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
				err := c.Find(&bson.M{"hotelid": hotelId, "inDate": indate, "outDate">
                                if err != nil {
                                        panic(err)
                                }
                                for _, r := range reserve {
                                        fmt.Printf("reservation check reservation number = %s\n", hot>                                        count += r.Number
                                }

                                // update memcached
                                err = s.MemcClient.Set(&memcache.Item{Key: memc_key, Value: []byte(st>
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
                        item, err = s.MemcClient.Get(memc_cap_key)
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
                                err = MemcClient.Set(&memcache.Item{Key: memc_cap_key, Value: []byte(>
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
                                res.HotelId = append(res.HotelId, hotelId)
                        }
                }
        }

        return res, nil
}

// MakeReservation makes a reservation based on given information
func MakeReservation(var req) (*pb.Result, error) {
        log.Println("MakeReservation")
        res := new(pb.Result)
        res.HotelId = make([]string, 0)

        // session, err := mgo.Dial("mongodb-reservation")
        // if err != nil {
        //      panic(err)
        // }
        // defer session.Close()
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
        hotelId := req.HotelId[0]

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
                        err := c.Find(&bson.M{"hotelid": hotelId, "inDate": indate, "outDate": outdat>
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
                        err = MemcClient.Set(&memcache.Item{Key: memc_cap_key, Value: []byte(strconv.>
                        if err != nil {
                                log.Warn("MMC error: ", err)
                        }
                } else {
                        fmt.Printf("Memmcached error = %s\n", err)
                        panic(err)
                }

                if count+int(req.RoomNumber) > hotel_cap {
			fmt.Printf("Not enough space left\n")
                        return res, nil
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

        res.HotelId = append(res.HotelId, hotelId)

        return res, nil
}

// loadUsers loads hotel users from database
func loadUsers(client *mongo.Client) map[string]string {

        coll := client.Database("user-db").Collection("user")

        filter := bson.D{}
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
        users_cached = loadUsers(MongoClient)
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
        collection := MongoClient.Database("user-db").Collection("user")

        // listAll(collection)

        // unmarshal json profiles
        var user User
        filter := bson.D{primitive.E{Key: "username", Value: username}}
        // filter := bson.M{{"username": username}}
        err := collection.FindOne(context.Background(), filter).Decode(&user)
        if err != nil {
                log.Println("Failed get user: ", err)
                return user, false
        }
        return user, true
}

// CheckUser returns whether the username and password are correct.
func CheckUser(var req) (*pb.Result, error) {
        res := new(pb.Result)

        fmt.Printf("CheckUser: %+v", req)

        sum := sha256.Sum256([]byte(req.Password))
        pass := fmt.Sprintf("%x", sum)

        use_cache := false

        if use_cache {
                password := lookupCache(req.Username)
                res.Correct = pass == password
        } else {
                user, _ := lookUpDB(req.Username)
                res.Correct = pass == user.Password
        }

        fmt.Printf(" >> pass: %t\n", res.Correct)

        return res, nil
}

// GetProfiles returns hotel profiles for requested IDs
func GetProfiles(var req) (*pb.Result, error) {
        // session, err := mgo.Dial("mongodb-profile")
        // if err != nil {
        //      panic(err)
        // }
        // defer session.Close()
        // fmt.Printf("In GetProfiles\n")

        // fmt.Printf("In GetProfiles after setting c\n")

        res := new(pb.Result)
        hotels := make([]*pb.Hotel, 0)

        // one hotel should only have one profile

        for _, i := range req.HotelIds {
                // first check memcached
                item, err := MemcClient.Get(i)
                if err == nil {
                        // memcached hit
                        // profile_str := string(item.Value)

                        // fmt.Printf("memc hit\n")
                        // fmt.Println(profile_str)

                        hotel_prof := new(pb.Hotel)
                        if err = json.Unmarshal(item.Value, hotel_prof); err != nil {
                                log.Warn(err)
                        }
                        hotels = append(hotels, hotel_prof)

                } else if err == memcache.ErrCacheMiss {
                        // memcached miss, set up mongo connection
                        session := MongoSession.Copy()
                        defer session.Close()
                        c := session.DB("profile-db").C("hotels")

                        hotel_prof := new(pb.Hotel)
                        err := c.Find(bson.M{"id": i}).One(&hotel_prof)

                        if err != nil {
                                log.Println("Hotel data not found in memcached: ", err)
                        }

                        // for _, h := range hotels {
			//      res.Hotels = append(res.Hotels, h)
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

        res.Hotels = hotels
        // fmt.Printf("In GetProfiles after getting resp\n")
        return res, nil
}

// GetRates gets rates for hotels for specific date range.
func GetRates(var req) (*pb.Result, error) {
        res := new(pb.Result)
        // session, err := mgo.Dial("mongodb-rate")
        // if err != nil {
        //      panic(err)
        // }
        // defer session.Close()

        ratePlans := make(RatePlans, 0)

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
                                        rate_p := new(pb.RatePlan)
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
        res.RatePlans = ratePlans

        return res, nil
}

// newGeoIndex returns a geo index with points loaded
func newGeoIndex(session *mgo.Session) *geoindex.ClusteringIndex {

        s := session.Copy()
        defer s.Close()
        c := s.DB("geo-db").C("geo")

        var points []*Point
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

        mgo.Session* MongoSession = new mgo.Session()
        return newGeoIndex(MongoSession).KNearest(
                center,
                maxSearchResults,
                geoindex.Km(maxSearchRadius), func(p geoindex.Point) bool {
                        return true
                },
        )
}

// Nearby returns all hotels within a given distance.
func Nearby(var req) (*pb.Result, error) {
        // fmt.Printf("In geo Nearby\n")

        var (
                points = getNearbyPoints(float64(req.Lat), float64(req.Lon))
                res    = &pb.Result{}
        )

        // fmt.Printf("geo after getNearbyPoints, len = %d\n", len(points))

        for _, p := range points {
                // fmt.Printf("In geo Nearby return hotelId = %s\n", p.Id())
                res.HotelIds = append(res.HotelIds, p.Id())
        }

        return res, nil
}



// Nearby returns ids of nearby hotels ordered by ranking algo
func SearchNearby(var req) (*pb.SearchResult, error) {
        // find nearby hotels
        fmt.Printf("in Search Nearby\n")

        fmt.Printf("nearby lat = %f\n", req.Lat)
        fmt.Printf("nearby lon = %f\n", req.Lon)

        nearby, err := Nearby({Lat: req.Lat, Lon: req.Lon})
        if err != nil {
                fmt.Printf("nearby error: %v", err)
                return nil, err
        }

        // var ids []string
        for _, hid := range nearby.HotelIds {
                fmt.Printf("get Nearby hotelId = %s\n", hid)
                // ids = append(ids, hid)
        }

        // find rates for hotels
        r := rate.Request{
                HotelIds: nearby.HotelIds,
                // HotelIds: []string{"2"},
                InDate:  req.InDate,
                OutDate: req.OutDate,
        }

        rates, err := GetRates(&r)
        if err != nil {
                fmt.Printf("rates error: %v", err)
                return nil, err
        }
        // TODO(hw): add simple ranking algo to order hotel ids:
	// * geo distance
        // * price (best discount?)
        // * reviews

        // build the response
        res := new(pb.SearchResult)
        for _, ratePlan := range rates.RatePlans {
                // fmt.Printf("get RatePlan HotelId = %s, Code = %s\n", ratePlan.HotelId, ratePlan.Code)
                res.HotelIds = append(res.HotelIds, ratePlan.HotelId)
        }
        return res, nil
}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	ret := ""
	requestUrl := ""
	if req.json.request == "search" {
		ret := SearchNearby(req.json)
	}
	else if req.json.request == "recommend" {
		ret := GetRecommendations(req.json)
	}
	else if req.json.request == "reserve" {
		if req.json.requestType == "check" {
			ret := CheckAvailability(req.json)
		}
		else if req.json.requestType == "make" {
			ret := MakeReservation(req.json)
		}
	}
	else if req.json.request == "user" {
		ret := CheckUser(req.json)
	}
	fmt.Fprintf(res, ret) // echo to caller
}

