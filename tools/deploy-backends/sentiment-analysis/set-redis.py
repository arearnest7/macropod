import json
import pprint
import csv
import os
import sys
import redis

def main():
    r = redis.Redis(host=sys.argv[1], port=6379, password=sys.argv[2])

    with open("review.csv", "r") as f:
        data = f.read()
        r.set("review.csv", data)
main()
