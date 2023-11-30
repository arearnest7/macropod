#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/video-analytics/video-streaming-d
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/video-analytics/video-recog-partial
