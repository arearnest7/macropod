#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
sudo k3s kubectl apply -f ../deploy-backends/hotel-app/database.yaml
sudo k3s kubectl apply -f ../deploy-backends/hotel-app/memcached.yaml
go run ../deploy-backends/hotel-app/hotel-app-backend.go
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-frontend
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-geo
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-profile
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-rate
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-recommend
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-reserve
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-search
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-user
