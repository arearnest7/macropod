#!/bin/bash
kubectl delete -f ../kn/yamls/sentiment-db-s.yaml
kubectl delete -f ../kn/yamls/sentiment-main-rcposc.yaml
kubectl delete -f ../kn/yamls/sentiment-product-sentiment-prs.yaml
kubectl delete -f ../kn/yamls/sentiment-service-sentiment-srs.yaml

