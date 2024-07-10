#!/bin/bash
kubectl delete -f ../kn/yamls/election-gateway.yaml
kubectl delete -f ../kn/yamls/election-get-results.yaml
kubectl delete -f ../kn/yamls/election-vote-enqueuer.yaml
kubectl delete -f ../kn/yamls/election-vote-processor.yaml
