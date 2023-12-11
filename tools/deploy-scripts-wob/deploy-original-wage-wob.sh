#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/wage-pay/wage-avg-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/wage-pay/wage-format-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/wage-pay/wage-merit-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/wage-pay/wage-stats-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/wage-pay/wage-sum-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/wage-pay/wage-validator-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/wage-pay/wage-write-merit-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/wage-pay/wage-write-raw-wob

