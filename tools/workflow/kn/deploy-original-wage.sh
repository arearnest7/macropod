#!/bin/bash
kubectl apply -f ./yamls/wage-avg.yaml 2> /dev/null
kubectl apply -f ./yamls/wage-format.yaml 2> /dev/null
kubectl apply -f ./yamls/wage-merit.yaml 2> /dev/null
kubectl apply -f ./yamls/wage-stats.yaml 2> /dev/null
kubectl apply -f ./yamls/wage-sum.yaml 2> /dev/null
kubectl apply -f ./yamls/wage-validator.yaml 2> /dev/null
kubectl apply -f ./yamls/wage-write-merit.yaml 2> /dev/null
kubectl apply -f ./yamls/wage-write-raw.yaml 2> /dev/null

