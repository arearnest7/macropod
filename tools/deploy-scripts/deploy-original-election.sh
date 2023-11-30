#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/serverless-election/election-gateway
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/serverless-election/election-get-results
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/serverless-election/election-vote-enqueuer
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/serverless-election/election-vote-processor
