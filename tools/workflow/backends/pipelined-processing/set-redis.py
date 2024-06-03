import json
import pprint
import csv
import os
import sys
import redis

def main():
    r = redis.Redis(host=sys.argv[1], port=6379, password=sys.argv[2])

    video = ["min_0.mp4", "min_0909.mp4", "min_1.mp4", "min_2.mp4", "min_3.mp4", "min_4.mp4", "min_5.mp4"]
    for v in video:
        with open(v, "rb") as f:
            data = f.read()
            r.set("original-" + v, data)
main()
