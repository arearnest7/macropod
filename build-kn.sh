#!/bin/bash
kn func build --push=true --path benchmarks/kn/original/video-analytics/video-streaming
kn func build --push=true --path benchmarks/kn/original/video-analytics/video-decoder
kn func build --push=true --path benchmarks/kn/original/video-analytics/video-recog
