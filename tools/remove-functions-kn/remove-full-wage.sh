#!/bin/bash
PATH=${1:-kn}
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/wage-full

