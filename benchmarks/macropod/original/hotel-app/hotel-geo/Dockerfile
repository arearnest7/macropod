FROM sysdevtamu/macropod-go:http
WORKDIR /app
COPY go.mod .
COPY func.go ./function
RUN go mod tidy
RUN go build -o main .
CMD ["/app/main"]
