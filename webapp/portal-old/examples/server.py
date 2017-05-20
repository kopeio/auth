#!/usr/bin/env python

import SimpleHTTPServer
import SocketServer

PORT = 8080

Handler = SimpleHTTPServer.SimpleHTTPRequestHandler
Handler.extensions_map.update({
        '': 'text/html',
        });

httpd = SocketServer.TCPServer(("", PORT), Handler)

print "serving at port", PORT
httpd.serve_forever()

