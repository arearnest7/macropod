#!/bin/bash
kubectl apply -f ./yamls/wage-stats-partial.yaml
kubectl apply -f ./yamls/wage-sum-amw.yaml
kubectl apply -f ./yamls/wage-validator-fw.yaml

