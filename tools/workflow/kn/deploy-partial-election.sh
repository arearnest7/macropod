#!/bin/bash
BUILD=${2:-false}
PUSH=${3:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/partial-reduced/serverless-election/election-gateway-vevp
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/partial-reduced/serverless-election/election-get-results-partial
