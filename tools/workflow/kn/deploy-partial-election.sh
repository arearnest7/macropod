#!/bin/bash
kubectl apply -f ./yamls/election-gateway-vevp.yaml
kubectl apply -f ./yamls/election-get-results-partial.yaml
