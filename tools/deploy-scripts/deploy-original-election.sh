#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/serverless-election/election-gateway
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/serverless-election/election-get-results
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/serverless-election/election-vote-enqueuer
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/serverless-election/election-vote-processor
