#!/bin/bash
kn func delete -p ../../../benchmarks/$1/partial-reduced/pipelined-processing/pipelined-checksum-partial
kn func delete -p ../../../benchmarks/$1/partial-reduced/pipelined-processing/pipelined-encrypt-partial
kn func delete -p ../../../benchmarks/$1/partial-reduced/pipelined-processing/pipelined-main-partial
kn func delete -p ../../../benchmarks/$1/partial-reduced/pipelined-processing/pipelined-zip-partial

