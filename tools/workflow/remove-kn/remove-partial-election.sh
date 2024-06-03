#!/bin/bash
kn func delete -p ../../../benchmarks/$1/partial-reduced/serverless-election/election-gateway-vevp
kn func delete -p ../../../benchmarks/$1/partial-reduced/serverless-election/election-get-results-partial
