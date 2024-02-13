package function

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
    "google.golang.org/grpc"

    pb "function/app_pb"
)

type Context struct {
    request string
    workflow_id string
    request_type string
    is_json bool
}

type server struct {
    pb.UnimplementedGRPCFunctionServer
}

func HTTPFunctionHandler(res http.ResponseWriter, req *http.Request) {
    logging_name, logging := os.LookupEnv("LOGGING_NAME")
    redisClient := redis.NewClient(&redis.Options{})
    c := context.Background()
    b, _ := ioutil.ReadAll(req.Body)
    body := string(b)
    if logging {
        logging_url := os.Getenv("LOGGING_URL")
        logging_password := os.Getenv("LOGGING_PASSWORD")
        redisClient = redis.NewClient(&redis.Options{
            Addr: logging_url,
            Password: logging_password,
            DB: 0,
        })
    }
    if req.Header.Get("Content-Type") == "application/json" {
        workflow_id := strconv.Itoa(rand.Intn(10000000))
        if logging {
            redisClient.Append(c, logging_name, "EXECUTION_HTTP_JSON_START," + workflow_id + ",NA,1," + time.Now().String() + "\n")
        }
        reply, _ := function_handler(Context{request: body, workflow_id: workflow_id, request_type: "HTTP", is_json: true})
        if logging {
            redisClient.Append(c, logging_name, "EXECUTION_HTTP_JSON_END," + workflow_id + ",NA,1," + time.Now().String() + "\n")
        }
        fmt.Fprintf(res, reply)
    } else {
        workflow_id := strconv.Itoa(rand.Intn(10000000))
        if logging {
            redisClient.Append(c, logging_name, "EXECUTION_HTTP_TEXT_START," + workflow_id + ",NA,1," + time.Now().String() + "\n")
        }
        reply, _ := function_handler(Context{request: body, workflow_id: workflow_id, request_type: "HTTP", is_json: false})
        if logging {
            redisClient.Append(c, logging_name, "EXECUTION_HTTP_TEXT_END," + workflow_id + ",NA,1," + time.Now().String() + "\n")
        }
        fmt.Fprintf(res, reply)
    }
}

func (s *server) GRPCFunctionHandler(ctx context.Context, in *pb.RequestBody) (*pb.ResponseBody, error) {
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
        redisClient.Append(c, logging_name, "EXECUTION_GRPC_START," + in.GetWorkflowId() + ",NA,1," + time.Now().String() + "\n")
    }
    reply, code := function_handler(Context{request: string(in.GetBody()), workflow_id: in.GetWorkflowId(), request_type: "GRPC", is_json: false})
    if logging {
        redisClient.Append(c, logging_name, "EXECUTION_GRPC_END," + in.GetWorkflowId() + ",NA,1," + time.Now().String() + "\n")
    }
    return &pb.ResponseBody{Reply: reply, Code: int32(code)}, nil
}

func main() {
    service_type, typed := os.LookupEnv("SERVICE_TYPE")
    if !typed || service_type == "HTTP" {
        http.HandleFunc("/", HTTPFunctionHandler)
        http.ListenAndServe(":" + os.Getenv("FUNC_PORT"), nil)
    } else if service_type == "GRPC" {
        l, _ := net.Listen("tcp", fmt.Sprintf(":%d", os.Getenv("FUNC_PORT")))
        s := grpc.NewServer()
        pb.RegisterGRPCFunctionServer(s, &server{})
        s.Serve(l)
    }
}
