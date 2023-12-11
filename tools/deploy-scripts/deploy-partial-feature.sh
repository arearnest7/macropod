#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/feature-generation/feature-extractor-partial
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/feature-generation/feature-orchestrator-wsr

