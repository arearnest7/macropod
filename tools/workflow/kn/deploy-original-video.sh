#!/bin/bash
kubectl apply -f ./yamls/video-streaming.yaml
kubectl apply -f ./yamls/video-decoder.yaml
kubectl apply -f ./yamls/video-recog.yaml
