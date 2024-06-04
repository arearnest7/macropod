#!/bin/bash
BUILD=${2:-false}
PUSH=${3:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/partial-reduced/feature-generation/feature-extractor-partial
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/partial-reduced/feature-generation/feature-orchestrator-wsr

