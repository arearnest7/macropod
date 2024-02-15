FROM arearnest7/faastroute-go:latest
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main .
CMD ["/app/main"]
