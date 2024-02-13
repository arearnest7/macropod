#!/bin/bash
PATH=${1:-kn}
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/video-streaming
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/video-decoder
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/video-recog
