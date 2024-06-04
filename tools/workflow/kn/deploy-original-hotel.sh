#!/bin/bash
BUILD=${2:-false}
PUSH=${3:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/original/hotel-app/hotel-frontend
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/original/hotel-app/hotel-geo
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/original/hotel-app/hotel-profile
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/original/hotel-app/hotel-rate
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/original/hotel-app/hotel-recommend
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/original/hotel-app/hotel-reserve
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/original/hotel-app/hotel-search
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/original/hotel-app/hotel-user
