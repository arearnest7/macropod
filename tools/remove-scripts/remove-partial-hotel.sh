#!/bin/bash
sudo k3s kubectl delete -f ../deploy-backends/hotel-app/database.yaml
sudo k3s kubectl delete -f ../deploy-backends/hotel-app/memcached.yaml
sudo kn func delete hotel-frontend-spgr
sudo kn func delete hotel-recommend-partial
sudo kn func delete hotel-reserve-partial
sudo kn func delete hotel-user-partial
