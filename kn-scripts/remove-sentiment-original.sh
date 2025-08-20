#bin/bash
sudo kubectl delete -f yamls/sentiment-cfail.yaml
sudo kubectl delete -f yamls/sentiment-db.yaml
sudo kubectl delete -f yamls/sentiment-main.yaml
sudo kubectl delete -f yamls/sentiment-product-or-service.yaml
sudo kubectl delete -f yamls/sentiment-product-result.yaml
sudo kubectl delete -f yamls/sentiment-product-sentiment.yaml
sudo kubectl delete -f yamls/sentiment-read-csv.yaml
sudo kubectl delete -f yamls/sentiment-service-result.yaml
sudo kubectl delete -f yamls/sentiment-service-sentiment.yaml
sudo kubectl delete -f yamls/sentiment-sfail.yaml
sudo kubectl delete -f yamls/sentiment-sns.yaml

