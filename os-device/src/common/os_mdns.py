import wifi
import mdns
import time

# wi-fi credentials
SSID = "ExampleSSID"
PASSWORD = "examplepassword"

# mdns
HOSTNAME = "opensesame-device"
INSTANCE_NAME = "OpenSesame Device"
PORT = 80

# overwrite defaults if settings exists
try:
    from os_settings import settings as _settings
    SSID = _settings.get("ssid", SSID)
    PASSWORD = _settings.get("password", PASSWORD)
    HOSTNAME = _settings.get("hostname", HOSTNAME)
    INSTANCE_NAME = _settings.get("instance_name", INSTANCE_NAME)
    PORT = _settings.get("port", PORT)
except Exception:
    # secrets.py won't exist unless running via the built output, so skip in dev
    pass

def connect_wifi(ssid, password, retries=3):
    # Reset radio
    wifi.radio.enabled = False
    time.sleep(0.5)
    wifi.radio.enabled = True

    for attempt in range(1, retries+1):
        try:
            print("Connecting to Wi‑Fi... (attempt {})".format(attempt))
            try:
                wifi.radio.disconnect()
            except Exception:
                pass
            wifi.radio.connect(ssid, password)
            print("Connected! IP address:", wifi.radio.ipv4_address)
            return
        except Exception as e:
            print("Wi‑Fi connection failed:", e)
            time.sleep(2)

    raise ConnectionError("Failed to connect to Wi‑Fi after {} retries".format(retries))

def start_mdns(hostname="opensesame-device", instance="OpenSesame Device", port=80):
    server = mdns.Server(wifi.radio)
    server.hostname = hostname
    server.instance_name = instance
    server.advertise_service(service_type="_http", protocol="_tcp", port=port)
    print(f"mDNS started: {hostname}.local → {instance} on port {port}")
    return server