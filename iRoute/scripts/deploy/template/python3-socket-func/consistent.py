import bisect
import hashlib

def get_hash(raw_str):
    md5_str = hashlib.md5(raw_str.encode()).hexdigest()
    return int(md5_str, 16)

class HashTable:
    def __init__(self):
        self.worker_list = []
        self.worker_table = {}
        self.virtual_num = 10

    def add_server(self, func_id):
        for index in range(0, self.virtual_num):
            worker_hash = get_hash("%s_%s" % (func_id, index))
            bisect.insort(self.worker_list, worker_hash)
            self.worker_table[worker_hash] = func_id
    
    def del_server(self, func_id):
        for index in range(0, self.virtual_num):
            worker_hash = get_hash("%s_%s" % (func_id, index))
            self.worker_list.remove(worker_hash)
            del self.worker_table[worker_hash]

    def get_server(self, source_key):
        key_hash = get_hash(source_key)
        index = bisect.bisect_left(self.worker_list, key_hash)
        index = index % len(self.worker_list)
        return self.worker_table[self.worker_list[index]]
