#!/bin/bash
kubectl delete -f ../kn/yamls/pipelined-checksum.yaml
kubectl delete -f ../kn/yamls/pipelined-encrypt.yaml
kubectl delete -f ../kn/yamls/pipelined-main.yaml
kubectl delete -f ../kn/yamls/pipelined-zip.yaml

