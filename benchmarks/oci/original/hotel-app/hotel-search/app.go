package function

import (
    "net/http"
    "os"
    "github.com/redis/go-redis/v9"
    "time"
    "context"
)

func function(res http.ResponseWriter, req *http.Request) {
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
    if logging {
        redisClient.Append(c, logging_name, time.Now().String() + "," + "0" + "," + "0" + "," + "0" + "," + "http" + "," + "0" + "\n")
    }
    FunctionHandler(res, req, "application/json", true)
    if logging {
        redisClient.Append(c, logging_name, time.Now().String() + "," + "0" + "," + "0" + "," + "0" + "," + "http" + "," + "1" + "\n")
    }
}

func main() {
    http.HandleFunc("/", function)
    http.ListenAndServe(":" + os.Getenv("PORT"), nil)
}
