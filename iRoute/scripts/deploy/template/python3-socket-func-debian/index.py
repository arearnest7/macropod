import threading
from http.server import HTTPServer, BaseHTTPRequestHandler
import threading

from function import handler

from config import Config
from sidecar import Sidecar

class HealthHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/_/health":
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b"OK")
        else:
            self.send_response(404)
            self.end_headers()

    def log_message(self, format, *args):
        pass

if __name__ == '__main__':
    health_server = HTTPServer(("0.0.0.0", 8080), HealthHandler)
    health_thread = threading.Thread(target=health_server.serve_forever, daemon=True)
    health_thread.start()

    config = Config()
    sidecar = Sidecar(handler.handle, config=config)
    sidecar.run()
