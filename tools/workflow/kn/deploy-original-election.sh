#!/bin/bash
sudo kubectl apply -f ./yamls/election-gateway.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/election-get-results.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/election-vote-enqueuer.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/election-vote-processor.yaml 2> /dev/null
