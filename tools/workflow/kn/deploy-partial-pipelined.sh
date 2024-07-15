#!/bin/bash
kubectl apply -f ./yamls/pipelined-checksum-partial.yaml 2> /dev/null
kubectl apply -f ./yamls/pipelined-encrypt-partial.yaml 2> /dev/null
kubectl apply -f ./yamls/pipelined-main-partial.yaml 2> /dev/null
kubectl apply -f ./yamls/pipelined-zip-partial.yaml 2> /dev/null

