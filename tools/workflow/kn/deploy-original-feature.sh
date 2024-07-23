#!/bin/bash
sudo kubectl apply -f ./yamls/feature-extractor.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/feature-orchestrator.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/feature-reducer.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/feature-status.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/feature-wait.yaml 2> /dev/null

