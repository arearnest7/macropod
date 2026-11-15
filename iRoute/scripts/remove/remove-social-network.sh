#!/bin/bash
faas-cli delete -f ../yamls/sn-entry.yml
faas-cli delete -f ../yamls/sn-cc.yml
faas-cli delete -f ../yamls/sn-funcs.yml
sudo rm -rf /mnt/sn
