package function

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	ret := ""
	requestUrl := ""
	if req.json.request == "search" {
		contents, err := ioutil.ReadFile("/etc/secret-volume/hotel-search")
        	if err != nil {
                	log.Fatal(err)
        	}
        	hotel-search := string(contents)
		requestUrl := hotel-search
	}
	else if req.json.request == "recommend" {
		contents, err := ioutil.ReadFile("/etc/secret-volume/hotel-recommend")
                if err != nil {
                        log.Fatal(err)
                }
                hotel-recommend := string(contents)
                requestUrl := hotel-recommend
	}
	else if req.json.request == "reserve" {
		contents, err := ioutil.ReadFile("/etc/secret-volume/hotel-reserve")
                if err != nil {
                        log.Fatal(err)
                }
                hotel-reserve := string(contents)
                requestUrl := hotel-reserve
	}
	else if req.json.request == "user" {
		contents, err := ioutil.ReadFile("/etc/secret-volume/hotel-user")
                if err != nil {
                        log.Fatal(err)
                }
                hotel-user := string(contents)
                requestUrl := hotel-user
	}

	req_url, err := http.NewRequest(http.MethodPost, requestURL, req.json)
        if err != nil {
                log.Fatal(err)
        }
        ret, err := http.DefaultClient.Do(req_url)
	fmt.Fprintf(res, ret) // echo to caller
}

