#!/bin/bash
kubectl apply -f ./yamls/hotel-frontend-spgr.yaml
kubectl apply -f ./yamls/hotel-recommend-partial.yaml
kubectl apply -f ./yamls/hotel-reserve-partial.yaml
kubectl apply -f ./yamls/hotel-user-partial.yaml
