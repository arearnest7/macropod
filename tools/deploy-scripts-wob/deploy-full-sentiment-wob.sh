#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/full-reduced/sentiment-analysis/sentiment-full-wob

