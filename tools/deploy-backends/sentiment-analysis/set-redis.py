import json
import pprint
import csv
import os
import sys
import redis

def main():
    r = redis.Redis(host='10.125.189.107', port=6379, password='redispassword1234')

    with open("review.csv", "r") as f:
        data = f.read()
        r.set("review.csv", data)
main()
