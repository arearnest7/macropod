package function

import (
    "os"
    "context"
    "time"
    "github.com/redis/go-redis/v9"
    "strconv"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    pb "function/app_pb"
)

func invoke(stub pb.GRPCFunctionClient, ctx context.Context, in *pb.RequestBody) (*pb.ResponseBody) {
    res, _ := stub.GRPCFunctionHandler(ctx, in)
    return res
}

func RPC(dest string, payloads []string, id string) ([]string) {
    logging_name, logging := os.LookupEnv("LOGGING_NAME")
    redisClient := redis.NewClient(&redis.Options{})
    c := context.Background()
    if logging {
        logging_url := os.Getenv("LOGGING_URL")
        logging_password := os.Getenv("LOGGING_PASSWORD")
        redisClient = redis.NewClient(&redis.Options{
            Addr: logging_url,
            Password: logging_password,
            DB: 0,
        })
    }
    channel, _ := grpc.Dial(dest, grpc.WithTransportCredentials(insecure.NewCredentials()))
    defer channel.Close()
    stub := pb.NewGRPCFunctionClient(channel)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    tl := make(chan *pb.ResponseBody)
    if logging {
        redisClient.Append(c, logging_name, "INVOKE_START," + id + "," + dest + "," + strconv.Itoa(len(payloads)) + "," + time.Now().String() + "\n")
    }
    for i := 0; i < len(payloads); i++ {
        go func(n int) {
            tl <- invoke(stub, ctx, &pb.RequestBody{Body: []byte(payloads[n]), WorkflowId: id})
        }(i)
    }
    results := make([]string, 0)
    for i := 0; i < len(payloads); i++ {
        res := <-tl
        results = append(results, res.GetReply())
    }
    if logging {
        redisClient.Append(c, logging_name, "INVOKE_END," + id + "," + dest + "," + strconv.Itoa(len(payloads)) + "," + time.Now().String() + "\n")
    }
    return results
}
