#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/partial-reduced/serverless-election/election-gateway-vevp-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/partial-reduced/serverless-election/election-get-results-partial-wob
