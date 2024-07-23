#!/bin/bash
sudo kubectl delete -f ../kn/yamls/feature-extractor.yaml
sudo kubectl delete -f ../kn/yamls/feature-orchestrator.yaml
sudo kubectl delete -f ../kn/yamls/feature-reducer.yaml
sudo kubectl delete -f ../kn/yamls/feature-status.yaml
sudo kubectl delete -f ../kn/yamls/feature-wait.yaml

