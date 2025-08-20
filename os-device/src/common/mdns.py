import time
import wifi
import mdns

# wi-fi credentials
SSID = "ExampleSSID"
PASSWORD = "examplepassword"

# mdns
HOSTNAME = "opensesame-device"
INSTANCE_NAME = "OpenSesame Device"
PORT = 80

# overwrite defaults if settings exists
try:
    from settings import settings as _settings
    SSID = _settings.get("ssid", SSID)
    PASSWORD = _settings.get("password", PASSWORD)
    HOSTNAME = _settings.get("hostname", HOSTNAME)
    INSTANCE_NAME = _settings.get("instance_name", INSTANCE_NAME)
    PORT = _settings.get("port", PORT)
except Exception:
    # secrets.py won't exist unless running via the built output, so skip in dev
    pass

def connect_wifi(ssid, password):
    wifi.radio.enabled = True
    print("Connecting to Wi‑Fi...")
    try:
        wifi.radio.connect(ssid, password)
    except Exception as e:
        print("Wi‑Fi connection failed:", e)
        raise
    print("Connected! IP address:", wifi.radio.ipv4_address)


def start_mdns(hostname="opensesame-device", instance="OpenSesame Device", port=80):
    server = mdns.Server(wifi.radio)
    server.hostname = hostname
    server.instance_name = instance
    server.advertise_service(service_type="_http", protocol="_tcp", port=port)
    print(f"mDNS started: {hostname}.local → {instance} on port {port}")
    return server