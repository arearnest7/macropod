#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/serverless-election/election-gateway-vevp
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/serverless-election/election-get-results-partial
