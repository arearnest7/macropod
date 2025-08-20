#bin/bash
sudo kubectl delete -f yamls/video-streaming.yaml
sudo kubectl delete -f yamls/video-decoder.yaml
sudo kubectl delete -f yamls/video-recog.yaml
