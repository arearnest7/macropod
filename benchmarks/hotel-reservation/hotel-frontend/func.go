package function

import (

)

func FunctionHandler(context Context) (string, int) {
    requestURL := ""
    body_u := context.JSON

    if body_u["Request"].(string) == "search" {
        requestURL = "HOTEL_SEARCH"
    } else if body_u["Request"].(string) == "recommend" {
        requestURL = "HOTEL_RECOMMEND"
    } else if body_u["Request"].(string) == "reserve" {
        requestURL = "HOTEL_RESERVE"
    } else if body_u["Request"].(string) == "user" {
        requestURL = "HOTEL_USER"
    } else if body_u["Request"].(string) == "profile" {
        requestURL = "HOTEL_PROFILE"
    }

    ret_val, _ := Invoke_JSON(context, requestURL, body_u)
    return ret_val, 200
}

