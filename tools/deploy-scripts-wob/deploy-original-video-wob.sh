#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/video-analytics/video-streaming-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/video-analytics/video-decoder-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/video-analytics/video-recog-wob
