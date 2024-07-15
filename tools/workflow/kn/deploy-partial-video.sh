#!/bin/bash
kubectl apply -f ./yamls/video-streaming-d.yaml 2> /dev/null
kubectl apply -f ./yamls/video-recog-partial.yaml 2> /dev/null
