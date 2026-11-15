#!/bin/bash
sudo kubectl delete -f ../yamls/excamera-entry.yaml
sudo kubectl delete -f ../yamls/excamera-group.yaml
sudo kubectl delete -f ../yamls/excamera-rebase.yaml
sudo kubectl delete -f ../yamls/excamera-reencode.yaml
sudo kubectl delete -f ../yamls/excamera-vpxenc.yaml
sudo kubectl delete -f ../yamls/excamera-vx-con.yaml
sudo kubectl delete -f ../yamls/excamera-xcdec.yaml
