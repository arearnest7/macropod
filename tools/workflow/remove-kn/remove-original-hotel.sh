#!/bin/bash
kubectl delete -f ../kn/yamls/hotel-frontend.yaml
kubectl delete -f ../kn/yamls/hotel-geo.yaml
kubectl delete -f ../kn/yamls/hotel-profile.yaml
kubectl delete -f ../kn/yamls/hotel-rate.yaml
kubectl delete -f ../kn/yamls/hotel-recommend.yaml
kubectl delete -f ../kn/yamls/hotel-reserve.yaml
kubectl delete -f ../kn/yamls/hotel-search.yaml
kubectl delete -f ../kn/yamls/hotel-user.yaml
