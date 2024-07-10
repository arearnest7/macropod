#!/bin/bash
kubectl apply -f ./yamls/hotel-frontend.yaml
kubectl apply -f ./yamls/hotel-geo.yaml
kubectl apply -f ./yamls/hotel-profile.yaml
kubectl apply -f ./yamls/hotel-rate.yaml
kubectl apply -f ./yamls/hotel-recommend.yaml
kubectl apply -f ./yamls/hotel-reserve.yaml
kubectl apply -f ./yamls/hotel-search.yaml
kubectl apply -f ./yamls/hotel-user.yaml
