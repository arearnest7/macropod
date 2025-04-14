FROM sysdevtamu/macropod-go:latest
COPY Dockerfile .
COPY go.mod .
COPY func.py ./function
RUN go mod tidy
RUN go build -o main .
