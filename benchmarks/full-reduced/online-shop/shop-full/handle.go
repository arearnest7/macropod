package function

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)



func AdServiceCall(var payload) {
	ret, err := exec.Command("./AdService " + payload).Output()
        if err != nil {
                fmt.Printf("%s", err)
        }
        return ret
}

func RecommendationServiceCall(var playload) {
	ret, err := exec.Command("python3 recommendation.py " + payload).Output()
        if err != nil {
                fmt.Printf("%s", err)
        }
        return ret
}

func ProductCatalogServiceCall(var payload) {
	ret, err := exec.Command("go run productcatalog.go " + payload).Output()
        if err != nil {
                fmt.Printf("%s", err)
        }
        return ret
}

func CartServiceCall(var payload) {
	ret, err := exec.Command("./CartService " + payload).Output()
        if err != nil {
                fmt.Printf("%s", err)
        }
        return ret
}

func ShippingServiceCall(var payload) {
	ret, err := exec.Command("go run shipping.go " + payload).Output()
        if err != nil {
                fmt.Printf("%s", err)
        }
        return ret
}

func CurrencyServiceCall(var payload) {
	ret, err := exec.Command("npm run currency.js " + payload).Output()
        if err != nil {
                fmt.Printf("%s", err)
        }
        return ret
}

func CheckoutServiceCall(var payload) {
	ret, err := exec.Command("go run checkout.go " + payload).Output()
        if err != nil {
                fmt.Printf("%s", err)
        }
        return ret
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

