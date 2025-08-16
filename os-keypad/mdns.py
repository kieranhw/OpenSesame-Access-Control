import time
import wifi
import mdns

wifi.radio.enabled = True

SSID = "ExampleSSID"
PASSWORD = "examplepassword"

print("Connecting to Wi‑Fi...")
try:
    wifi.radio.connect(SSID, PASSWORD)
except Exception as e:
    print("Wi‑Fi connection failed:", e)
    while True:
        time.sleep(1)

print("Connected! IP address:", wifi.radio.ipv4_address)

# Set up mDNS.
mdns_server = mdns.Server(wifi.radio)
mdns_server.hostname = "opensesame-keypad"
mdns_server.instance_name = "Opensesame Keypad"
mdns_server.advertise_service(service_type="_http", protocol="_tcp", port=80)

print("mDNS service started. You should be able to reach it at opensesame-keypad.local")

while True:
    time.sleep(1)