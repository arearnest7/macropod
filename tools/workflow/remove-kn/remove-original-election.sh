#!/bin/bash
sudo kubectl delete -f ../kn/yamls/election-gateway.yaml
sudo kubectl delete -f ../kn/yamls/election-get-results.yaml
sudo kubectl delete -f ../kn/yamls/election-vote-enqueuer.yaml
sudo kubectl delete -f ../kn/yamls/election-vote-processor.yaml
