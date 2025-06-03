FROM sysdevtamu/macropod-go:unified
COPY Dockerfile .
COPY go.mod .
COPY func.go ./function
RUN go mod tidy 2> /dev/null
RUN go build -o main .
