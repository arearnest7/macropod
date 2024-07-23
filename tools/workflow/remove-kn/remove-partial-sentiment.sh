#!/bin/bash
sudo kubectl delete -f ../kn/yamls/sentiment-db-s.yaml
sudo kubectl delete -f ../kn/yamls/sentiment-main-rcposc.yaml
sudo kubectl delete -f ../kn/yamls/sentiment-product-sentiment-prs.yaml
sudo kubectl delete -f ../kn/yamls/sentiment-service-sentiment-srs.yaml

