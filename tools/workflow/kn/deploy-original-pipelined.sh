#!/bin/bash
sudo kubectl apply -f ./yamls/pipelined-checksum.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/pipelined-encrypt.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/pipelined-main.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/pipelined-zip.yaml 2> /dev/null

