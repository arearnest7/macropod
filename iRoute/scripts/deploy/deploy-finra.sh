#!/bin/bash
sudo mkdir -p /mnt/finra
sudo chown 100:101 /mnt/finra
export TARGET_NODE=$(hostname)
faas-cli deploy -f ../yamls/finra-cc.yml
faas-cli deploy -f ../yamls/finra-funcs.yml
faas-cli deploy -f ../yamls/finra-entry.yml
