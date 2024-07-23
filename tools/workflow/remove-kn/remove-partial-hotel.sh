#!/bin/bash
sudo kubectl delete -f ../kn/yamls/hotel-frontend-spgr.yaml
sudo kubectl delete -f ../kn/yamls/hotel-recommend-partial.yaml
sudo kubectl delete -f ../kn/yamls/hotel-reserve-partial.yaml
sudo kubectl delete -f ../kn/yamls/hotel-user-partial.yaml
