#!/bin/bash
sudo kubectl delete -f ../kn/yamls/hotel-frontend.yaml
sudo kubectl delete -f ../kn/yamls/hotel-geo.yaml
sudo kubectl delete -f ../kn/yamls/hotel-profile.yaml
sudo kubectl delete -f ../kn/yamls/hotel-rate.yaml
sudo kubectl delete -f ../kn/yamls/hotel-recommend.yaml
sudo kubectl delete -f ../kn/yamls/hotel-reserve.yaml
sudo kubectl delete -f ../kn/yamls/hotel-search.yaml
sudo kubectl delete -f ../kn/yamls/hotel-user.yaml
