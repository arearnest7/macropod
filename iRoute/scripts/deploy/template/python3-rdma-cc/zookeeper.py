from kazoo.client import KazooClient
from kazoo.security import make_digest_acl

class ZookeeperClient:
    def __init__(self, hosts, username, password):
        self.zk = KazooClient(hosts=hosts)
        self.zk.start()
        self.zk.add_auth('digest', f'{username}:{password}')

    def make_write_acl(self, username, password):
        """
        Create a write ACL for the given username and password.
        """
        return make_digest_acl(username, password, all=True)
    
    def make_read_acl(self, username, password):
        """
        Create a read ACL for the given username and password.
        """
        return make_digest_acl(username, password, read=True)

    def create_node(self, path, data=None, acl=[None]):
        if self.zk.exists(path):
            self.zk.delete(path, recursive=True)
        if data is None:
            data = b''
        else:
            data = data.encode('utf-8')
        self.zk.create(path, data, acl=acl, makepath=True)

    def get_node(self, path):
        if self.zk.exists(path):
            data, stat = self.zk.get(path)
            return data.decode('utf-8')
        else:
            return None
        
    def set_node(self, path, data):
        if self.zk.exists(path):
            self.zk.set(path, data.encode('utf-8'))
        else:
            raise Exception(f"Node {path} does not exist")
        
    def delete_node(self, path):
        if self.zk.exists(path):
            self.zk.delete(path, recursive=True)
        else:
            raise Exception(f"Node {path} does not exist")
        
    def stop(self):
        """
        Stop the Zookeeper client.
        """
        self.zk.stop()
