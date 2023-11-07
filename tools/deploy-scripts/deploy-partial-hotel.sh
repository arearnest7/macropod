#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
sudo k3s kubectl apply -f ../deploy-backends/hotel-app/database.yaml
sudo k3s kubectl apply -f ../deploy-backends/hotel-app/memcached.yaml
go run ../deploy-backends/hotel-app/hotel-app-backend.go
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/hotel-app/hotel-frontend-spgr
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/hotel-app/hotel-recommend-partial
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/hotel-app/hotel-reserve-partial
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/hotel-app/hotel-user-partial
