#!/bin/bash
PATH=${1:-kn}
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/election-gateway-vevp
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/election-get-results-partial
