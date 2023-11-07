#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/serverless-election/election-gateway-vevp
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/serverless-election/election-get-results-partial
