#!/bin/bash
kubectl apply -f ./yamls/hotel-frontend-spgr.yaml 2> /dev/null
kubectl apply -f ./yamls/hotel-recommend-partial.yaml 2> /dev/null
kubectl apply -f ./yamls/hotel-reserve-partial.yaml 2> /dev/null
kubectl apply -f ./yamls/hotel-user-partial.yaml 2> /dev/null
