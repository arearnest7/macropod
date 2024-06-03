import sys
import redis

def main():
    r = redis.Redis(host=sys.argv[1], port=6379, password=sys.argv[2])

    amazon = ["reviews10mb.csv", "reviews20mb.csv", "reviews50mb.csv", "reviews100mb.csv"]
    for a in amazon:
        with open(a, "r") as f:
            data = f.read()
            r.set("amazon-" + a, data)
main()
