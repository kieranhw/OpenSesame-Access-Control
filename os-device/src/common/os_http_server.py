import wifi

class HTTPServer:
    def __init__(self, pool, port=80):
        self.server = pool.socket(pool.AF_INET, pool.SOCK_STREAM)
        self.server.settimeout(None)
        self.server.bind(("0.0.0.0", port))
        self.server.listen(1)
        self.port = port
        self.ip = str(wifi.radio.ipv4_address)
        print("HTTP server running at http://{}:{}".format(self.ip, self.port))

    def send_response(self, conn, body, status="200 OK", content_type="text/plain"):
        response = (
            "HTTP/1.1 {}\r\n"
            "Content-Type: {}\r\n"
            "Content-Length: {}\r\n"
            "Connection: close\r\n"
            "\r\n"
            "{}"
        ).format(status, content_type, len(body), body)
        try:
            conn.send(response.encode("utf-8"))
        except Exception as e:
            print("Send error:", e)
        try:
            conn.close()
        except Exception:
            pass

    def serve_forever(self, handler):
        try:
            while True:
                conn, addr = self.server.accept()
                print("Connection from", addr)

                try:
                    buffer = bytearray(1024)
                    bytes_read = conn.recv_into(buffer)
                    if bytes_read == 0:
                        conn.close()
                        continue
                    request = buffer[:bytes_read].decode("utf-8")
                except Exception as e:
                    print("Error reading request:", e)
                    conn.close()
                    continue

                if not request:
                    conn.close()
                    continue

                request_line = request.split("\r\n")[0]
                print("Request line:", request_line)

                try:
                    method, path, _ = request_line.split(" ")
                except ValueError:
                    conn.close()
                    continue

                # Extract body if present
                body = ""
                if "\r\n\r\n" in request:
                    body = request.split("\r\n\r\n", 1)[1]

                # Call user handler
                try:
                    status, content_type, response_body = handler(method, path, body)
                except Exception as e:
                    print("Handler error:", e)
                    status, content_type, response_body = (
                        "500 Internal Server Error",
                        "text/plain",
                        "Server error: {}".format(e),
                    )

                self.send_response(conn, response_body, status=status, content_type=content_type)

        finally:
            self.close()

    def close(self):
        print("Closing HTTP server...")
        try:
            self.server.close()
        except Exception as e:
            print("Error closing server:", e)