#!/bin/bash
sudo kubectl delete -f ../yamls/finra-entry.yaml
sudo kubectl delete -f ../yamls/finra-lastpx.yaml
sudo kubectl delete -f ../yamls/finra-margin-balance.yaml
sudo kubectl delete -f ../yamls/finra-marketdata.yaml
sudo kubectl delete -f ../yamls/finra-side.yaml
sudo kubectl delete -f ../yamls/finra-trddate.yaml
sudo kubectl delete -f ../yamls/finra-volume.yaml
