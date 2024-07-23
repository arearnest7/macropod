#!/bin/bash
sudo kubectl apply -f ./yamls/feature-extractor-partial.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/feature-orchestrator-wsr.yaml 2> /dev/null

