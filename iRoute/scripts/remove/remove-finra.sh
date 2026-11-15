#!/bin/bash
faas-cli delete -f ../yamls/finra-entry.yml
faas-cli delete -f ../yamls/finra-cc.yml
faas-cli delete -f ../yamls/finra-funcs.yml
sudo rm -rf /mnt/finra
