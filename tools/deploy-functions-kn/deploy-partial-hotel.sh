#!/bin/bash
PATH=${1:-kn}
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$PATH/partial-reduced/hotel-app/hotel-frontend-spgr
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$PATH/partial-reduced/hotel-app/hotel-recommend-partial
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$PATH/partial-reduced/hotel-app/hotel-reserve-partial
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$PATH/partial-reduced/hotel-app/hotel-user-partial
