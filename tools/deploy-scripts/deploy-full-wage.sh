#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/full-reduced/wage-pay/wage-full

