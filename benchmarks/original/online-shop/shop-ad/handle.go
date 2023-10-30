package function

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"os/exec"
)

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	ret, err := exec.Command("./AdService " + req.json).Output()
	if err != nil {
        	fmt.Printf("%s", err)
    	}
	fmt.Fprintf(res, ret) // echo to caller
}
