#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/pipelined-processing/pipelined-checksum
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/pipelined-processing/pipelined-encrypt
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/pipelined-processing/pipelined-main
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/pipelined-processing/pipelined-zip

