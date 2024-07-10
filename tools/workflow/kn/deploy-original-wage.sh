#!/bin/bash
kubectl apply -f ./yamls/wage-avg.yaml
kubectl apply -f ./yamls/wage-format.yaml
kubectl apply -f ./yamls/wage-merit.yaml
kubectl apply -f ./yamls/wage-stats.yaml
kubectl apply -f ./yamls/wage-sum.yaml
kubectl apply -f ./yamls/wage-validator.yaml
kubectl apply -f ./yamls/wage-write-merit.yaml
kubectl apply -f ./yamls/wage-write-raw.yaml

