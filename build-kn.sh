#!/bin/bash
kn func build --push=true --path benchmarks/kn/unified/video-analytics/video-unified

kn func build --push=true --path benchmarks/kn/original/video-analytics/video-streaming
kn func build --push=true --path benchmarks/kn/original/video-analytics/video-decoder
kn func build --push=true --path benchmarks/kn/original/video-analytics/video-recog
