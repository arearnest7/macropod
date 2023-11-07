#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/video-analytics/video-streaming
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/video-analytics/video-decoder
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/original/video-analytics/video-recog
