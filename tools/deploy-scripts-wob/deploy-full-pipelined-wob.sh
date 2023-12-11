#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/full-reduced/pipelined-processing/pipelined-full-wob

