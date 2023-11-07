#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
sudo k3s kubectl apply -f ../deploy-backends/hotel-app/database.yaml
sudo k3s kubectl apply -f ../deploy-backends/hotel-app/memcached.yaml
go run ../deploy-backends/hotel-app/hotel-app-backend.go
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/full-reduced/hotel-app/hotel-full
