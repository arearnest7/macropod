#!/bin/bash
kubectl delete -f ../kn/yamls/wage-avg.yaml
kubectl delete -f ../kn/yamls/wage-format.yaml
kubectl delete -f ../kn/yamls/wage-merit.yaml
kubectl delete -f ../kn/yamls/wage-stats.yaml
kubectl delete -f ../kn/yamls/wage-sum.yaml
kubectl delete -f ../kn/yamls/wage-validator.yaml
kubectl delete -f ../kn/yamls/wage-write-merit.yaml
kubectl delete -f ../kn/yamls/wage-write-raw.yaml

