#!/bin/bash
sudo mkdir -p /mnt/sn
sudo chown 100:101 /mnt/sn
export TARGET_NODE=$(hostname)
faas-cli deploy -f ../yamls/sn-cc.yml
faas-cli deploy -f ../yamls/sn-funcs.yml
faas-cli deploy -f ../yamls/sn-entry.yml
