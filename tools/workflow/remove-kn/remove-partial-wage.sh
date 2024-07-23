#!/bin/bash
sudo kubectl delete -f ../kn/yamls/wage-stats-partial.yaml
sudo kubectl delete -f ../kn/yamls/wage-sum-amw.yaml
sudo kubectl delete -f ../kn/yamls/wage-validator-fw.yaml

