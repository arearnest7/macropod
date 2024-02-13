#!/bin/bash
PATH=${1:-kn}
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/election-gateway
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/election-get-results
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/election-vote-enqueuer
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/election-vote-processor
