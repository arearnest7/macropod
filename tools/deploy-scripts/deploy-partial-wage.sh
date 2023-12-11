#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/wage-pay/wage-stats-partial
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/wage-pay/wage-sum-amw
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/wage-pay/wage-validator-fw

