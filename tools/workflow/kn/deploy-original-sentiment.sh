#!/bin/bash
sudo kubectl apply -f ./yamls/sentiment-cfail.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/sentiment-db.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/sentiment-main.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/sentiment-product-or-service.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/sentiment-product-result.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/sentiment-product-sentiment.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/sentiment-read-csv.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/sentiment-service-result.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/sentiment-service-sentiment.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/sentiment-sfail.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/sentiment-sns.yaml 2> /dev/null

