#bin/bash
host=${1:-127.0.0.1}
for i in yamls/sentiment*; do sed -i "s/knative-functions\.127\.0\.0\.1/knative-functions\.$host/g" $i; done;

sudo kubectl apply -f yamls/sentiment-cfail.yaml 2>/dev/null
sudo kubectl apply -f yamls/sentiment-db.yaml 2>/dev/null
sudo kubectl apply -f yamls/sentiment-main.yaml 2>/dev/null
sudo kubectl apply -f yamls/sentiment-product-or-service.yaml 2>/dev/null
sudo kubectl apply -f yamls/sentiment-product-result.yaml 2>/dev/null
sudo kubectl apply -f yamls/sentiment-product-sentiment.yaml 2>/dev/null
sudo kubectl apply -f yamls/sentiment-read-csv.yaml 2>/dev/null
sudo kubectl apply -f yamls/sentiment-service-result.yaml 2>/dev/null
sudo kubectl apply -f yamls/sentiment-service-sentiment.yaml 2>/dev/null
sudo kubectl apply -f yamls/sentiment-sfail.yaml 2>/dev/null
sudo kubectl apply -f yamls/sentiment-sns.yaml 2>/dev/null

for i in yamls/sentiment*; do sed -i "s/knative-functions\.$host/knative-functions\.127\.0\.0\.1/g" $i; done;
