#!/bin/bash
kubectl apply -f ./yamls/feature-extractor.yaml
kubectl apply -f ./yamls/feature-orchestrator.yaml
kubectl apply -f ./yamls/feature-reducer.yaml
kubectl apply -f ./yamls/feature-status.yaml
kubectl apply -f ./yamls/feature-wait.yaml

