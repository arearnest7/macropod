#!/bin/bash
kubectl delete -f ../kn/yamls/hotel-frontend-spgr.yaml
kubectl delete -f ../kn/yamls/hotel-recommend-partial.yaml
kubectl delete -f ../kn/yamls/hotel-reserve-partial.yaml
kubectl delete -f ../kn/yamls/hotel-user-partial.yaml
