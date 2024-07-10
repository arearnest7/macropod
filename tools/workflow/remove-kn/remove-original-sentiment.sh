#!/bin/bash
kubectl delete -f ../kn/yamls/sentiment-cfail.yaml
kubectl delete -f ../kn/yamls/sentiment-db.yaml
kubectl delete -f ../kn/yamls/sentiment-main.yaml
kubectl delete -f ../kn/yamls/sentiment-product-or-service.yaml
kubectl delete -f ../kn/yamls/sentiment-product-result.yaml
kubectl delete -f ../kn/yamls/sentiment-product-sentiment.yaml
kubectl delete -f ../kn/yamls/sentiment-read-csv.yaml
kubectl delete -f ../kn/yamls/sentiment-service-result.yaml
kubectl delete -f ../kn/yamls/sentiment-service-sentiment.yaml
kubectl delete -f ../kn/yamls/sentiment-sfail.yaml
kubectl delete -f ../kn/yamls/sentiment-sns.yaml

