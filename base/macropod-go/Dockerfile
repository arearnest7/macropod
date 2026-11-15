FROM golang:latest
WORKDIR /app
COPY . .
RUN go mod tidy 2> /dev/null
RUN go build -o main
CMD ["/app/main"]
