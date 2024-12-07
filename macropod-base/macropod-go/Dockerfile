FROM golang:alpine
WORKDIR /app
COPY . .
RUN go get math/rand
RUN go mod tidy
RUN go build -o main
CMD ["/app/main"]
