FROM golang:1.24
WORKDIR /app
COPY . .
RUN apt-get update
RUN apt-get install curl -y
RUN go mod tidy
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \ 
    go build -o main
CMD ["/app/main"]
