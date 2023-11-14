#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/video-analytics/video-streaming-d
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/video-analytics/video-recog-partial
