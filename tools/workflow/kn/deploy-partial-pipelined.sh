#!/bin/bash
kubectl apply -f ./yamls/pipelined-checksum-partial.yaml
kubectl apply -f ./yamls/pipelined-encrypt-partial.yaml
kubectl apply -f ./yamls/pipelined-main-partial.yaml
kubectl apply -f ./yamls/pipelined-zip-partial.yaml

