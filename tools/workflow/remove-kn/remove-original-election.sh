#!/bin/bash
kn func delete -p ../../benchmarks/$1/original/serverless-election/election-gateway
kn func delete -p ../../benchmarks/$1/original/serverless-election/election-get-results
kn func delete -p ../../benchmarks/$1/original/serverless-election/election-vote-enqueuer
kn func delete -p ../../benchmarks/$1/original/serverless-election/election-vote-processor
