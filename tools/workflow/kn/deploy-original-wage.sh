#!/bin/bash
BUILD=${2:-false}
PUSH=${3:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/original/wage-pay/wage-avg
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/original/wage-pay/wage-format
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/original/wage-pay/wage-merit
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/original/wage-pay/wage-stats
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/original/wage-pay/wage-sum
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/original/wage-pay/wage-validator
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/original/wage-pay/wage-write-merit
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/original/wage-pay/wage-write-raw

