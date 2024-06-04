#!/bin/bash
BUILD=${2:-false}
PUSH=${3:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/partial-reduced/wage-pay/wage-stats-partial
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/partial-reduced/wage-pay/wage-sum-amw
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/partial-reduced/wage-pay/wage-validator-fw

