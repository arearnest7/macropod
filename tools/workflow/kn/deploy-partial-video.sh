#!/bin/bash
sudo kubectl apply -f ./yamls/video-streaming-d.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/video-recog-partial.yaml 2> /dev/null
