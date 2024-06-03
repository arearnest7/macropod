#!/bin/bash
sudo k3s kubectl apply -f database.yaml
sudo k3s kubectl apply -f memcached.yaml
