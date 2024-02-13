FROM golang:alpine
WORKDIR /app
COPY . .
RUN go get github.com/redis/go-redis/v9
RUN go get google.golang.org/grpc
RUN go get math/rand
RUN go build -o main . 
CMD ["/app/main"]
