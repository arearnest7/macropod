#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-frontend
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-geo
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-profile
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-rate
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-recommend
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-reserve
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-search
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-user
