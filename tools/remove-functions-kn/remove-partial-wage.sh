#!/bin/bash
PATH=${1:-kn}
kn func delete -f ../../benchmarks/$PATH/full-reduced/serverless-election/wage-stats-partial
kn func delete -f ../../benchmarks/$PATH/full-reduced/serverless-election/wage-sum-amw
kn func delete -f ../../benchmarks/$PATH/full-reduced/serverless-election/wage-validator-fw

