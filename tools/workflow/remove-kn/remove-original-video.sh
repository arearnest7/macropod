#!/bin/bash
kn func delete -p ../../../benchmarks/$1/original/video-analytics/video-streaming
kn func delete -p ../../../benchmarks/$1/original/video-analytics/video-decoder
kn func delete -p ../../../benchmarks/$1/original/video-analytics/video-recog
