FROM golang:alpine
WORKDIR /app
COPY . .
RUN go mod tidy
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \ 
    go build -o main
CMD ["/app/main"]