package function

import (
    "fmt"
    "encoding/json"
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

// Nearby returns ids of nearby hotels ordered by ranking algo
func Nearby(context Context) string {
    req := context.JSON
    // find nearby hotels
    fmt.Printf("in Search Nearby\n")

    //fmt.Printf("nearby lat = %f\n", req.Lat)
    //fmt.Printf("nearby lon = %f\n", req.Lon)

    requestURL := "HOTEL_GEO"
    payload := map[string]interface{}{"Lat": req["Lat"].(float64), "Lon": req["Lon"].(float64)}
    //req_url, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(body_g))
    //req_url.Header.Add("Content-Type", "application/json")
    //client := &http.Client{}
    //nearby, err := client.Do(req_url)
    //if err != nil {
    //        fmt.Printf("nearby error: %v", err)
    //        return ""
    //}
    nearby, _ := Invoke_JSON(context, requestURL, payload)

    // var ids []string
    //nearbyBody, err := ioutil.ReadAll(nearby.Body)
    nearby_u := make([]string, 0)
    json.Unmarshal([]byte(nearby), &nearby_u)
    for _, hid := range nearby_u {
        fmt.Printf("get Nearby hotelId = %s\n", hid)
        // ids = append(ids, hid)
    }

    // find rates for hotels
    r := map[string]interface{}{
        "HotelIds": nearby_u,
        // HotelIds: []string{"2"},
        "InDate":  req["InDate"].(string),
        "OutDate": req["OutDate"].(string),
    }

    requestURL2 := "HOTEL_RATE"
    //req_url2, err := http.NewRequest(http.MethodPost, requestURL2, bytes.NewBuffer(body_r))
    //req_url2.Header.Add("Content-Type", "application/json")
    //ratesRet, err := client.Do(req_url2)
    //if err != nil {
    //        fmt.Printf("rates error: %v", err)
    //        return ""
    //}
    fmt.Println(string(r))
    ratesRet, _ := Invoke_JSON(context, requestURL2, r)
    //rates, _ := ioutil.ReadAll(ratesRet)
    // TODO(hw): add simple ranking algo to order hotel ids:
    // * geo distance
    // * price (best discount?)
    // * reviews

    // build the response
    res := make([]string, 0)
    var rate_p RatePlans
    json.Unmarshal([]byte(ratesRet), &rate_p)
    for _, ratePlan := range rate_p {
        // fmt.Printf("get RatePlan HotelId = %s, Code = %s\n", ratePlan.HotelId, ratePlan.Code)
        res = append(res, ratePlan.hotelId)
    }
    ret, _ := json.Marshal(res)
    return string(ret)
}

func FunctionHandler(context Context) (string, int) {
    return Nearby(context), 200
}
