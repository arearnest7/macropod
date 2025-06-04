FROM golang:1.24
WORKDIR /app
COPY . .
RUN go mod tidy 2> /dev/null
RUN go build -o main
CMD ["/app/main"]
