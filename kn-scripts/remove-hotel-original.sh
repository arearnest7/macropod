#bin/bash
sudo kubectl delete -f yamls/hotel-frontend.yaml
sudo kubectl delete -f yamls/hotel-geo.yaml
sudo kubectl delete -f yamls/hotel-profile.yaml
sudo kubectl delete -f yamls/hotel-rate.yaml
sudo kubectl delete -f yamls/hotel-recommend.yaml
sudo kubectl delete -f yamls/hotel-reserve.yaml
sudo kubectl delete -f yamls/hotel-search.yaml
sudo kubectl delete -f yamls/hotel-user.yaml
