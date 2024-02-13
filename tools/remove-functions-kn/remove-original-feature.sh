#!/bin/bash
PATH=${1:-kn}
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/feature-extractor
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/feature-orchestrator
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/feature-reducer
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/feature-status
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/feature-wait

