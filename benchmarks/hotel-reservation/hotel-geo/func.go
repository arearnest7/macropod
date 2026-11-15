package function

import (
    "fmt"
    "os"
    "io/ioutil"
    "encoding/json"

    log "github.com/sirupsen/logrus"

    "github.com/hailocab/go-geoindex"
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

// Implement Point interface
func (p *Point) Lat() float64 { return p.Plat }
func (p *Point) Lon() float64 { return p.Plon }
func (p *Point) Id() string   { return p.Pid }

// newGeoIndex returns a geo index with points loaded
func newGeoIndex() *geoindex.ClusteringIndex {

    //s := session.Copy()
    //defer s.Close()
    f, _ := os.Open("geo_db.json")
    c, _ := ioutil.ReadAll(f)

    points := make([]*Point, 0)
    err := json.Unmarshal(c, &points)
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

    //MongoSession, _ := mgo.Dial(os.Getenv("HOTEL_APP_DATABASE"))
    return newGeoIndex().KNearest(
        center,
        maxSearchResults,
        geoindex.Km(maxSearchRadius), func(p geoindex.Point) bool {
            return true
        },
    )
}

// Nearby returns all hotels within a given distance.
func Nearby(req map[string]interface{}) string {
    // fmt.Printf("In geo Nearby\n")

    var (
        points = getNearbyPoints(req["Lat"].(float64), req["Lon"].(float64))
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

func FunctionHandler(context Context) (string, int) {
    return Nearby(context.JSON), 200
}
