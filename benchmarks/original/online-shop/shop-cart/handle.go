package function

import (
        "context"
        "fmt"
        "net/http"
        "strings"
        "os/exec"
)

type CartItem struct {
	product_id string
	quantity int
}

type Cart struct {
	user_id string
	items []CartItem
}

func AddItem(var req) string {

}

func GetCart(var req) string {

}

func EmptyCart(var req) string {

}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
        ret = ""
	if req.json.request == "add" {
		ret = AddItem(req.json)
	}
	else if req.json.request == "cart" {
		ret = GetCart(req.json)
	}
	else if req.json.request == "empty" {
		ret = EmptyCart(req.json)
	}

	fmt.Fprintf(res, ret) // echo to caller
}
