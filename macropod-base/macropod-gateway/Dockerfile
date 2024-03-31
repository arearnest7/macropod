FROM golang:alpine
WORKDIR /app
COPY . .
RUN go get github.com/redis/go-redis/v9 google.golang.org/grpc math/rand github.com/go-mmap/mmap
RUN go mod tidy
RUN go build -o main
CMD ["/app/main"]
