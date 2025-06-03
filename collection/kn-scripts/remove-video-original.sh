#bin/bash
kubectl delete -f yamls/video-streaming.yaml
kubectl delete -f yamls/video-decoder.yaml
kubectl delete -f yamls/video-recog.yaml
