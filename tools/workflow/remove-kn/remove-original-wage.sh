#!/bin/bash
sudo kubectl delete -f ../kn/yamls/wage-avg.yaml
sudo kubectl delete -f ../kn/yamls/wage-format.yaml
sudo kubectl delete -f ../kn/yamls/wage-merit.yaml
sudo kubectl delete -f ../kn/yamls/wage-stats.yaml
sudo kubectl delete -f ../kn/yamls/wage-sum.yaml
sudo kubectl delete -f ../kn/yamls/wage-validator.yaml
sudo kubectl delete -f ../kn/yamls/wage-write-merit.yaml
sudo kubectl delete -f ../kn/yamls/wage-write-raw.yaml

