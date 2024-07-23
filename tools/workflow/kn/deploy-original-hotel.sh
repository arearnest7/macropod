#!/bin/bash
sudo kubectl apply -f ./yamls/hotel-frontend.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/hotel-geo.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/hotel-profile.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/hotel-rate.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/hotel-recommend.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/hotel-reserve.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/hotel-search.yaml 2> /dev/null
sudo kubectl apply -f ./yamls/hotel-user.yaml 2> /dev/null
