#!/bin/bash
PATH=${1:-kn}
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$PATH/original/video-analytics/video-streaming
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$PATH/original/video-analytics/video-decoder
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$PATH/original/video-analytics/video-recog
