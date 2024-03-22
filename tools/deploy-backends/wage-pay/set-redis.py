import json
import pprint
import csv
import os
import sys
import redis

def main():
    r = redis.Redis(host=sys.argv[1], port=6379, password=sys.argv[2])

    for i in range(100):
        with open(str(i), "r") as f:
            data = f.read()
            r.set("raw-" + str(i), data)
main()
