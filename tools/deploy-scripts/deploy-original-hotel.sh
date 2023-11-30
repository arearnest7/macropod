#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-frontend
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-geo
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-profile
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-rate
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-recommend
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-reserve
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-search
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/hotel-app/hotel-user
