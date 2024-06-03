#!/bin/bash
kn func delete -p ../../../benchmarks/$1/original/pipelined-processing/pipelined-checksum
kn func delete -p ../../../benchmarks/$1/original/pipelined-processing/pipelined-encrypt
kn func delete -p ../../../benchmarks/$1/original/pipelined-processing/pipelined-main
kn func delete -p ../../../benchmarks/$1/original/pipelined-processing/pipelined-zip

