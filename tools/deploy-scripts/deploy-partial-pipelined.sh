#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/pipelined-processing/pipelined-checksum-partial
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/pipelined-processing/pipelined-encrypt-partial
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/pipelined-processing/pipelined-main-partial
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/pipelined-processing/pipelined-zip-partial

