#bin/bash
kubectl delete -f yamls/sentiment-cfail.yaml
kubectl delete -f yamls/sentiment-db.yaml
kubectl delete -f yamls/sentiment-main.yaml
kubectl delete -f yamls/sentiment-product-or-service.yaml
kubectl delete -f yamls/sentiment-product-result.yaml
kubectl delete -f yamls/sentiment-product-sentiment.yaml
kubectl delete -f yamls/sentiment-read-csv.yaml
kubectl delete -f yamls/sentiment-service-result.yaml
kubectl delete -f yamls/sentiment-service-sentiment.yaml
kubectl delete -f yamls/sentiment-sfail.yaml
kubectl delete -f yamls/sentiment-sns.yaml

