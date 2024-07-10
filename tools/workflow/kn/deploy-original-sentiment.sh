#!/bin/bash
kubectl apply -f ./yamls/sentiment-cfail.yaml
kubectl apply -f ./yamls/sentiment-db.yaml
kubectl apply -f ./yamls/sentiment-main.yaml
kubectl apply -f ./yamls/sentiment-product-or-service.yaml
kubectl apply -f ./yamls/sentiment-product-result.yaml
kubectl apply -f ./yamls/sentiment-product-sentiment.yaml
kubectl apply -f ./yamls/sentiment-read-csv.yaml
kubectl apply -f ./yamls/sentiment-service-result.yaml
kubectl apply -f ./yamls/sentiment-service-sentiment.yaml
kubectl apply -f ./yamls/sentiment-sfail.yaml
kubectl apply -f ./yamls/sentiment-sns.yaml

