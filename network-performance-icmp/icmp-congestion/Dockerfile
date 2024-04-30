FROM golang:1.22
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go mod tidy
RUN go build -o main . 
USER root 
CMD ["/app/main"]

