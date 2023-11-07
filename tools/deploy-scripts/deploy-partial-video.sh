#!/bin/bash
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/video-analytics/video-streaming-d
sudo kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/partial-reduced/video-analytics/video-recog-partial
