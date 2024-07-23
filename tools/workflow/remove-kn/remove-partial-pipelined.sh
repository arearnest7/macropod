#!/bin/bash
sudo kubectl delete -f ../kn/yamls/pipelined-checksum-partial.yaml
sudo kubectl delete -f ../kn/yamls/pipelined-encrypt-partial.yaml
sudo kubectl delete -f ../kn/yamls/pipelined-main-partial.yaml
sudo kubectl delete -f ../kn/yamls/pipelined-zip-partial.yaml

