#!/bin/bash
sudo kubectl apply -f ./yamls/video-streaming.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/video-decoder.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/video-recog.yaml 2> /dev/null
