package function

import (
	"os"
	"encoding/json"
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

func FunctionHandler(context Context) (string, int) {
	requestURL := ""
	//body, _ := ioutil.ReadAll(req.Body)
	body_u := RequestBody{}
	json.Unmarshal([]byte(context.Request), &body_u)
  	//defer req.Body.Close()
	if body_u.Request == "search" {
		requestURL = os.Getenv("HOTEL_SEARCH")
	} else if body_u.Request == "recommend" {
                requestURL = os.Getenv("HOTEL_RECOMMEND")
	} else if body_u.Request == "reserve" {
                requestURL = os.Getenv("HOTEL_RESERVE")
	} else if body_u.Request == "user" {
                requestURL = os.Getenv("HOTEL_USER")
	}

	body_m, _ := json.Marshal(body_u)
	//req_url, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(body_m))
        //if err != nil {
        //        log.Fatal(err)
        //}
	//req_url.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
        //ret, err := client.Do(req_url)
        //retBody, err := ioutil.ReadAll(ret.Body)
        //ret_val, err := json.Marshal(retBody)
        ret_val := RPC(context, requestURL, [][]byte{body_m})[0]
        return ret_val, 200
}

