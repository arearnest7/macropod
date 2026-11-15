#!/bin/bash
faas-cli delete -f ../yamls/excamera-entry.yml
faas-cli delete -f ../yamls/excamera-cc.yml
faas-cli delete -f ../yamls/excamera-funcs.yml
sudo rm -rf /mnt/excamera
