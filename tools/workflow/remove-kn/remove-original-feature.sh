#!/bin/bash
kubectl delete -f ../kn/yamls/feature-extractor.yaml
kubectl delete -f ../kn/yamls/feature-orchestrator.yaml
kubectl delete -f ../kn/yamls/feature-reducer.yaml
kubectl delete -f ../kn/yamls/feature-status.yaml
kubectl delete -f ../kn/yamls/feature-wait.yaml

