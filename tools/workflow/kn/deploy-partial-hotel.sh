#!/bin/bash
BUILD=${2:-false}
PUSH=${3:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/partial-reduced/hotel-app/hotel-frontend-spgr
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/partial-reduced/hotel-app/hotel-recommend-partial
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/partial-reduced/hotel-app/hotel-reserve-partial
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/partial-reduced/hotel-app/hotel-user-partial
