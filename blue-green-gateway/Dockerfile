FROM golang:1.21
RUN apt-get update && \
    apt-get install -y vim && \
    rm -rf /var/lib/apt/lists/*
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go mod tidy
RUN go build -o main . 
USER root 
CMD ["/app/main"]

