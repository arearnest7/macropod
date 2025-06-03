#bin/bash
kubectl delete -f yamls/hotel-frontend.yaml
kubectl delete -f yamls/hotel-geo.yaml
kubectl delete -f yamls/hotel-profile.yaml
kubectl delete -f yamls/hotel-rate.yaml
kubectl delete -f yamls/hotel-recommend.yaml
kubectl delete -f yamls/hotel-reserve.yaml
kubectl delete -f yamls/hotel-search.yaml
kubectl delete -f yamls/hotel-user.yaml
