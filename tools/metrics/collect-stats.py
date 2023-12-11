import json
import sys
import os
import psutil
import time
import random

def worker_handler(file_name):
    try:
        with open(file_name, "w") as f:
            f.write("timestamp,user,idle,iowait,steal,guest,current,min,max,cpu_load_1,cpu_load_5,cpu_load_15,total,available,used,free,active,inactive,bytes_sent,bytes_recv,packets_sent,packets_recv\n")
            while True:
                ct = psutil.cpu_times()
                cf = psutil.cpu_freq()
                la = psutil.getloadavg()
                vm = psutil.virtual_memory()
                ic = psutil.net_io_counters()
                timestamp = time.strftime("%b %d %Y %H:%M:%S", time.localtime())
                f.write("%s,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f\n" % (timestamp, ct[0], ct[3], ct[4], ct[7], ct[8], cf[0], cf[1], cf[2], la[0], la[1], la[2], vm[0], vm[1], vm[3], vm[4], vm[5], vm[6], ic[0], ic[1], ic[2], ic[3]))
                time.sleep(0.005)
    except KeyboardInterrupt:
        pass

def main():
    if len(sys.argv) != 2:
        print("python collect-stats.py [file_name]")
    worker_handler(sys.argv[1])
main()
