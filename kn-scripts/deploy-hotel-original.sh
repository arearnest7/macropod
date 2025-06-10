#bin/bash
host=${1:-127.0.0.1}
for i in yamls/hotel*; do sed -i "s/knative-functions\.127\.0\.0\.1/knative-functions\.$host/g" $i; done;

kubectl apply -f yamls/hotel-frontend.yaml 2>/dev/null
kubectl apply -f yamls/hotel-geo.yaml 2>/dev/null
kubectl apply -f yamls/hotel-profile.yaml 2>/dev/null
kubectl apply -f yamls/hotel-rate.yaml 2>/dev/null
kubectl apply -f yamls/hotel-recommend.yaml 2>/dev/null
kubectl apply -f yamls/hotel-reserve.yaml 2>/dev/null
kubectl apply -f yamls/hotel-search.yaml 2>/dev/null
kubectl apply -f yamls/hotel-user.yaml 2>/dev/null

for i in yamls/hotel*; do sed -i "s/knative-functions\.$host/knative-functions\.127\.0\.0\.1/g" $i; done;
