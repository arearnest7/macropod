package function

import (
    "net/http"
    "os"
)

func function(res http.ResponseWriter, req *http.Request) {
    FunctionHandler(res, req, "application/json", true)
}

func main() {
    http.HandleFunc("/", function)
    http.ListenAndServe(":" + os.Getenv("PORT"), nil)
}
