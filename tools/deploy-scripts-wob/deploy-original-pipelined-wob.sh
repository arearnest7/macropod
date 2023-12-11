#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/pipelined-processing/pipelined-checksum-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/pipelined-processing/pipelined-encrypt-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/pipelined-processing/pipelined-main-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/pipelined-processing/pipelined-zip-wob

