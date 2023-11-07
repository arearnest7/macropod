#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
sudo k3s kubectl apply -f ../deploy-backends/online-shop/database.yaml
