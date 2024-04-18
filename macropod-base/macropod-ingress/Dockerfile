FROM golang:alpine
WORKDIR /app
COPY . .
RUN go get math/rand google.golang.org/grpc k8s.io/client-go@latest
RUN go mod tidy
RUN go build -o main
CMD ["/app/main"]
