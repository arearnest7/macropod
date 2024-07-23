#!/bin/bash
sudo kubectl delete -f ../kn/yamls/pipelined-checksum.yaml
sudo kubectl delete -f ../kn/yamls/pipelined-encrypt.yaml
sudo kubectl delete -f ../kn/yamls/pipelined-main.yaml
sudo kubectl delete -f ../kn/yamls/pipelined-zip.yaml

