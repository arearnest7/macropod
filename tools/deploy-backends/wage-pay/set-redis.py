import json
import pprint
import csv
import os
import sys
import redis

def main():
    r = redis.Redis(host='10.125.189.107', port=6379, password='redispassword1234')

    for i in range(100):
        with open(str(i), "r") as f:
            data = f.read()
            r.set("raw-" + str(i), data)
main()
