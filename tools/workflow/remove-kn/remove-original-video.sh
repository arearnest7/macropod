#!/bin/bash
kubectl delete -f ../kn/yamls/video-streaming.yaml
kubectl delete -f ../kn/yamls/video-decoder.yaml
kubectl delete -f ../kn/yamls/video-recog.yaml
