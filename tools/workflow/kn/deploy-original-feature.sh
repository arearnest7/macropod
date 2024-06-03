#!/bin/bash
BUILD=${2:-false}
PUSH=${3:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/original/feature-generation/feature-extractor
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/original/feature-generation/feature-orchestrator
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/original/feature-generation/feature-reducer
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/original/feature-generation/feature-status
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/original/feature-generation/feature-wait

