#!/bin/bash
PATH=${1:-kn}
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/pipelined-checksum
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/pipelined-encrypt
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/pipelined-main
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/pipelined-zip

