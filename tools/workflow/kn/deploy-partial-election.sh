#!/bin/bash
kubectl apply -f ./yamls/election-gateway-vevp.yaml 2> /dev/null
kubectl apply -f ./yamls/election-get-results-partial.yaml 2> /dev/null
