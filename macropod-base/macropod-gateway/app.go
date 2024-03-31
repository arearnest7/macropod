package main

import (
    "net/http"
    "os"
    "context"
    "fmt"
    "net"
    "time"
    "math/rand"
    "github.com/redis/go-redis/v9"
    "strconv"
    "io/ioutil"
    "github.com/go-mmap/mmap"
    "google.golang.org/grpc"

    pb "app/app_pb"
    function "app/function"
)

type server struct {
    pb.UnimplementedGRPCFunctionServer
}

func invoke_workflow() {

}

func deploy_workflow() {

}

func update_workflow() {

}

func delete_workflow() {

}

func main() {

}
