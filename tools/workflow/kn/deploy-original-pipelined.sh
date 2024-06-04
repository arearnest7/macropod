#!/bin/bash
BUILD=${2:-false}
PUSH=${3:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/original/pipelined-processing/pipelined-checksum
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/original/pipelined-processing/pipelined-encrypt
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/original/pipelined-processing/pipelined-main
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/original/pipelined-processing/pipelined-zip

