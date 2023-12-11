#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/partial-reduced/pipelined-processing/pipelined-checksum-partial-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/partial-reduced/pipelined-processing/pipelined-encrypt-partial-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/partial-reduced/pipelined-processing/pipelined-main-partial-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/partial-reduced/pipelined-processing/pipelined-zip-partial-wob

