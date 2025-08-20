import board
import digitalio
import json
import time
import wifi
import socketpool
from os_settings import settings
from os_gpio import gpio
from os_mdns import connect_wifi, start_mdns
from os_http_server import HTTPServer

def setup_output(pin_number):
    pin = digitalio.DigitalInOut(getattr(board, "GP{}".format(pin_number)))
    pin.direction = digitalio.Direction.OUTPUT
    pin.value = False # default to low
    return pin

red_led = setup_output(gpio["output"]["red_led"])
yellow_led = setup_output(gpio["output"]["yellow_led"])
green_led = setup_output(gpio["output"]["green_led"])
lock_relay = setup_output(gpio["output"]["lock_relay"])

connect_wifi(settings["ssid"], settings["password"])
start_mdns(
    hostname=settings["hostname"],
    instance=settings["instance_name"],
    port=settings["port"]
)
pool = socketpool.SocketPool(wifi.radio)

lock_state = {"state": "unknown"}

def set_locked():
    yellow_led.value = False
    red_led.value = False
    green_led.value = True
    
    lock_state["state"] = "locked"
    lock_relay.value = True

def set_unlocked():
    red_led.value = False
    green_led.value = False
    yellow_led.value = True
    
    lock_state["state"] = "unlocked"
    lock_relay.value = False

def flash_failed():
    for _ in range(10):  # 10 flashes
        red_led.value = True
        time.sleep(0.1)
        red_led.value = False
        time.sleep(0.1)
    red_led.value = False
    lock_state["state"] = "failed"

def request_handler(method, path, body):
    if method == "POST" and path == "/lock":
        set_locked()
        return ("200 OK", "application/json", json.dumps(lock_state))

    elif method == "POST" and path == "/unlock":
        set_unlocked()
        return ("200 OK", "application/json", json.dumps(lock_state))

    elif method == "POST" and path == "/command_failed":
        flash_failed()
        return ("200 OK", "application/json", json.dumps(lock_state))

    elif method == "GET" and path == "/status":
        return ("200 OK", "application/json", json.dumps(lock_state))

    else:
        return ("404 Not Found", "text/plain", "Not found")

server = HTTPServer(pool, settings["port"])
try:
    server.serve_forever(request_handler)
finally:
    print("Main shutting down...")
    try:
        server.close()
    except Exception:
        pass
    try:
        wifi.radio.disconnect()
    except Exception:
        pass
    wifi.radio.enabled = False