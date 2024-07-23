#!/bin/bash
sudo kubectl apply -f ./yamls/wage-stats-partial.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/wage-sum-amw.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/wage-validator-fw.yaml 2> /dev/null

