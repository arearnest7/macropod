#!/bin/bash
sudo kubectl apply -f ./yamls/pipelined-checksum-partial.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/pipelined-encrypt-partial.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/pipelined-main-partial.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/pipelined-zip-partial.yaml 2> /dev/null

