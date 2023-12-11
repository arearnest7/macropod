#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/hotel-app/hotel-frontend-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/hotel-app/hotel-geo-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/hotel-app/hotel-profile-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/hotel-app/hotel-rate-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/hotel-app/hotel-recommend-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/hotel-app/hotel-reserve-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/hotel-app/hotel-search-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/hotel-app/hotel-user-wob
