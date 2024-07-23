#!/bin/bash
sudo kubectl delete -f ../kn/yamls/sentiment-cfail.yaml
sudo kubectl delete -f ../kn/yamls/sentiment-db.yaml
sudo kubectl delete -f ../kn/yamls/sentiment-main.yaml
sudo kubectl delete -f ../kn/yamls/sentiment-product-or-service.yaml
sudo kubectl delete -f ../kn/yamls/sentiment-product-result.yaml
sudo kubectl delete -f ../kn/yamls/sentiment-product-sentiment.yaml
sudo kubectl delete -f ../kn/yamls/sentiment-read-csv.yaml
sudo kubectl delete -f ../kn/yamls/sentiment-service-result.yaml
sudo kubectl delete -f ../kn/yamls/sentiment-service-sentiment.yaml
sudo kubectl delete -f ../kn/yamls/sentiment-sfail.yaml
sudo kubectl delete -f ../kn/yamls/sentiment-sns.yaml

