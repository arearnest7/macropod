#!/bin/bash
sudo k3s kubectl delete -f ../deploy-backends/hotel-app/database.yaml
sudo k3s kubectl delete -f ../deploy-backends/hotel-app/memcached.yaml
sudo kn func delete hotel-frontend
sudo kn func delete hotel-geo
sudo kn func delete hotel-profile
sudo kn func delete hotel-rate
sudo kn func delete hotel-recommend
sudo kn func delete hotel-reserve
sudo kn func delete hotel-search
sudo kn func delete hotel-user
