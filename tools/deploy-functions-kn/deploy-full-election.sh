#!/bin/bash
PATH=${1:-kn}
BUILD=${2:-false}
PUSH=${3:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$PATH/full-reduced/serverless-election/election-full

