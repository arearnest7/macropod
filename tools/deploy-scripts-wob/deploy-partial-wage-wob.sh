#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/partial-reduced/wage-pay/wage-stats-partial-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/partial-reduced/wage-pay/wage-sum-amw-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/partial-reduced/wage-pay/wage-validator-fw-wob

