#!/bin/bash
PATH=${1:-kn}
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/pipelined-checksum-partial
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/pipelined-encrypt-partial
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/pipelined-main-partial
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/pipelined-zip-partial

