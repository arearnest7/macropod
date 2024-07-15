#!/bin/bash
kubectl apply -f ./yamls/feature-extractor-partial.yaml 2> /dev/null
kubectl apply -f ./yamls/feature-orchestrator-wsr.yaml 2> /dev/null

