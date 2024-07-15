#!/bin/bash
kubectl apply -f ./yamls/wage-stats-partial.yaml 2> /dev/null
kubectl apply -f ./yamls/wage-sum-amw.yaml 2> /dev/null
kubectl apply -f ./yamls/wage-validator-fw.yaml 2> /dev/null

