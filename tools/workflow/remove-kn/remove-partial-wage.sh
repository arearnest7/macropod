#!/bin/bash
kubectl delete -f ../kn/yamls/wage-stats-partial.yaml
kubectl delete -f ../kn/yamls/wage-sum-amw.yaml
kubectl delete -f ../kn/yamls/wage-validator-fw.yaml

