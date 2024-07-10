#!/bin/bash
kubectl apply -f ./yamls/sentiment-db-s.yaml
kubectl apply -f ./yamls/sentiment-main-rcposc.yaml
kubectl apply -f ./yamls/sentiment-product-sentiment-prs.yaml
kubectl apply -f ./yamls/sentiment-service-sentiment-srs.yaml

