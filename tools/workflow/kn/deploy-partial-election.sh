#!/bin/bash
sudo kubectl apply -f ./yamls/election-gateway-vevp.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/election-get-results-partial.yaml 2> /dev/null
