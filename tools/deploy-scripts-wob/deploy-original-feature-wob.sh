#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/feature-generation/feature-extractor-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/feature-generation/feature-orchestrator-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/feature-generation/feature-reducer-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/feature-generation/feature-status-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/feature-generation/feature-wait-wob

