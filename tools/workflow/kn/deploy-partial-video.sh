#!/bin/bash
BUILD=${2:-false}
PUSH=${3:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/partial-reduced/video-analytics/video-streaming-d
kn func deploy --build=$BUILD --push=$PUSH --path ../../../benchmarks/$1/partial-reduced/video-analytics/video-recog-partial
