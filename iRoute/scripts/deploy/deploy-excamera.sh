#!/bin/bash
sudo mkdir -p /mnt/excamera
sudo chown 100:101 /mnt/excamera
export TARGET_NODE=$(hostname)
faas-cli deploy -f ../yamls/excamera-cc.yml
faas-cli deploy -f ../yamls/excamera-funcs.yml
faas-cli deploy -f ../yamls/excamera-entry.yml
