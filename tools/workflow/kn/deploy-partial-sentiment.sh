#!/bin/bash
sudo kubectl apply -f ./yamls/sentiment-db-s.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/sentiment-main-rcposc.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/sentiment-product-sentiment-prs.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/sentiment-service-sentiment-srs.yaml 2> /dev/null

