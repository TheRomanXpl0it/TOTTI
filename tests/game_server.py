import http.server
import random
import socketserver
import json

PORT = 8000


class ServerHandler(http.server.BaseHTTPRequestHandler):
	def do_PUT(self):
		try:
			#print(self.headers)
			content_length = int(self.headers['Content-Length'])
			recived = self.rfile.read(content_length).decode('utf-8')
			#print(f"Recived: {recived}")
			recived = json.loads(recived)

			body = []
			for flag in recived:
				'''
				subAccepted:     "accepted",
				subInvalid:      "invalid",
				subOld:          "too old",
				subYourOwn:      "your own",
				subStolen:       "already claimed",
				subNop:          "from NOP team",
				subNotAvailable: "is not available",
				'''
				choice = random.choices(
					population=['accepted', 'invalid', 'too old', 'your own', 'already claimed', 'from NOP team'],
					weights=(1, 0.1, 0.1, 0.1, 0.1, 0.1),
					k=1
				)[0]
				body.append({"msg": f"[{flag}] {choice}", "flag": flag, "status": True})
			
			if random.random() < 0.1:
				body = {"code": "RATE_LIMIT", "message": "[RATE_LIMIT] Rate limit exceeded"}
				body = json.dumps(body)
				self.send_response(500)
				self.send_header("Content-type", "application/json; charset=utf-8")
				self.send_header("Content-Length", str(len(body)))
				self.end_headers()
				self.wfile.write(bytes(body, 'utf-8'))
				return

			body = json.dumps(body)
			print(body)
			self.send_response(200)
			self.send_header("Content-type", "application/json; charset=utf-8")
			self.send_header("Content-Length", str(len(body)))
			self.end_headers()
			self.wfile.write(bytes(body, 'utf-8'))

		except Exception as e:
			print(e)


handler = ServerHandler

httpd = socketserver.TCPServer(("", PORT), handler)

print("Serving at port", PORT)
httpd.serve_forever()