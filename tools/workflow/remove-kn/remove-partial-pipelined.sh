#!/bin/bash
kubectl delete -f ../kn/yamls/pipelined-checksum-partial.yaml
kubectl delete -f ../kn/yamls/pipelined-encrypt-partial.yaml
kubectl delete -f ../kn/yamls/pipelined-main-partial.yaml
kubectl delete -f ../kn/yamls/pipelined-zip-partial.yaml

