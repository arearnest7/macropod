FROM arearnest7/faastroute-go:latest
WORKDIR /app
COPY . .
RUN go build -o main .
CMD ["/app/main"]
