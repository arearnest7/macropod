#!/bin/bash
sudo kubectl delete -f ../kn/yamls/video-streaming.yaml
sudo kubectl delete -f ../kn/yamls/video-decoder.yaml
sudo kubectl delete -f ../kn/yamls/video-recog.yaml
