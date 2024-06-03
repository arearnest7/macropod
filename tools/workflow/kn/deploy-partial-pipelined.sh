#!/bin/bash
BUILD=${2:-false}
PUSH=${3:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/partial-reduced/pipelined-processing/pipelined-checksum-partial
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/partial-reduced/pipelined-processing/pipelined-encrypt-partial
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/partial-reduced/pipelined-processing/pipelined-main-partial
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/partial-reduced/pipelined-processing/pipelined-zip-partial

