package main

import (
    "net/http"
    "time"
    "fmt"
    "os"
)

func function(res http.ResponseWriter, req *http.Request) {
    fmt.Println(time.Now().String() + "," + "0" + "," + "0" + "," + "0" + "," + "HTTP" + "," + "0" + "\n")
    FunctionHandler(res, req, "application/json", true)
    fmt.Println(time.Now().String() + "," + "0" + "," + "0" + "," + "0" + "," + "HTTP" + "," + "1" + "\n")
}

func main() {
    http.HandleFunc("/", function)
    http.ListenAndServe(":" + os.Getenv("PORT"), nil)
}
