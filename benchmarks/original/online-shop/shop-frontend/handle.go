package function

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func AdServiceCall(var payload) {
	contents, err := ioutil.ReadFile("/etc/secret-volume/shop-ad")
        if err != nil {
                log.Fatal(err)
        }
        ad := string(contents)

        requestURL := ad + ":80"
        req_url, err := http.NewRequest(http.MethodPost, requestURL, payload)
        if err != nil {
                log.Fatal(err)
        }
        ret := http.DefaultClient.Do(req_url)
}

func RecommendationServiceCall(var payload) {
	contents, err := ioutil.ReadFile("/etc/secret-volume/shop-recommendation")
        if err != nil {
                log.Fatal(err)
        }
        recommendation := string(contents)

        requestURL := recommendation + ":80"
        req_url, err := http.NewRequest(http.MethodPost, requestURL, payload)
        if err != nil {
                log.Fatal(err)
        }
        ret := http.DefaultClient.Do(req_url)
}

func ProductCatalogServiceCall(var payload) {
	contents, err := ioutil.ReadFile("/etc/secret-volume/shop-productcatalog")
        if err != nil {
                log.Fatal(err)
        }
        productcatalog := string(contents)

        requestURL := productcatalog + ":80"
        req_url, err := http.NewRequest(http.MethodPost, requestURL, payload)
        if err != nil {
                log.Fatal(err)
        }
        ret := http.DefaultClient.Do(req_url)
}

func CartServiceCall(var payload) {
	contents, err := ioutil.ReadFile("/etc/secret-volume/shop-cart")
        if err != nil {
                log.Fatal(err)
        }
        cart := string(contents)

        requestURL := cart + ":80"
        req_url, err := http.NewRequest(http.MethodPost, requestURL, payload)
        if err != nil {
                log.Fatal(err)
        }
        ret := http.DefaultClient.Do(req_url)
}

func ShippingServiceCall(var payload) {
	contents, err := ioutil.ReadFile("/etc/secret-volume/shop-shipping")
        if err != nil {
                log.Fatal(err)
        }
        shipping := string(contents)

        requestURL := shipping + ":80"
        req_url, err := http.NewRequest(http.MethodPost, requestURL, payload)
        if err != nil {
                log.Fatal(err)
        }
        ret := http.DefaultClient.Do(req_url)
}

func CurrencyServiceCall(var payload) {
	contents, err := ioutil.ReadFile("/etc/secret-volume/shop-currency")
        if err != nil {
                log.Fatal(err)
        }
        currency := string(contents)

        requestURL := currency + ":80"
        req_url, err := http.NewRequest(http.MethodPost, requestURL, payload)
        if err != nil {
                log.Fatal(err)
        }
        ret := http.DefaultClient.Do(req_url)
}

func CheckoutServiceCall(var payload) {
	contents, err := ioutil.ReadFile("/etc/secret-volume/shop-checkout")
        if err != nil {
                log.Fatal(err)
        }
        checkout := string(contents)

        requestURL := checkout + ":80"
        req_url, err := http.NewRequest(http.MethodPost, requestURL, payload)
        if err != nil {
                log.Fatal(err)
        }
        ret := http.DefaultClient.Do(req_url)
}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	let ret := ""
	let request = req.json.request
	if request == "ad" {
		ret = AdServiceCall(req.json)
	}
	else if request == "recommendation" {
		ret = RecommendationServiceCall(req.json)
	}
	else if request == "productcatalog" {
		ret = ProductCatalogServiceCall(req.json)
	}
	else if request == "cart" {
		ret = CartServiceCall(req.json)
	}
	else if request == "shipping" {
		ret = ShippingServiceCall(req.json)
	}
	else if request == "currency" {
		ret = CurrencyServiceCall(req.json)
	}
	else if request == "checkout" {
		ret = CheckoutServiceCall(req.json)
	}
	fmt.Fprintf(res, ret) // echo to caller
}

