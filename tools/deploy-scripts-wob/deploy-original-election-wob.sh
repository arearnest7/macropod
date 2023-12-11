#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/serverless-election/election-gateway-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/serverless-election/election-get-results-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/serverless-election/election-vote-enqueuer-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/serverless-election/election-vote-processor-wob
