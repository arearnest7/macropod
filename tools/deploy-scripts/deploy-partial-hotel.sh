#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/hotel-app/hotel-frontend-spgr
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/hotel-app/hotel-recommend-partial
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/hotel-app/hotel-reserve-partial
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/hotel-app/hotel-user-partial
