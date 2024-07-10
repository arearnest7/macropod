#!/bin/bash
kubectl apply -f ./yamls/election-gateway.yaml
kubectl apply -f ./yamls/election-get-results.yaml
kubectl apply -f ./yamls/election-vote-enqueuer.yaml
kubectl apply -f ./yamls/election-vote-processor.yaml
