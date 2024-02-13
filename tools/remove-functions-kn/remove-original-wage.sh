#!/bin/bash
PATH=${1:-kn}
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/wage-avg
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/wage-format
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/wage-merit
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/wage-stats
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/wage-sum
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/wage-validator
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/wage-write-merit
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/wage-write-raw

