#!/bin/bash
kubectl apply -f ./yamls/pipelined-checksum.yaml
kubectl apply -f ./yamls/pipelined-encrypt.yaml
kubectl apply -f ./yamls/pipelined-main.yaml
kubectl apply -f ./yamls/pipelined-zip.yaml

