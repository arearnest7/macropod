#!/bin/bash
BUILD=${2:-false}
PUSH=${3:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/original/serverless-election/election-gateway
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/original/serverless-election/election-get-results
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/original/serverless-election/election-vote-enqueuer
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/original/serverless-election/election-vote-processor
