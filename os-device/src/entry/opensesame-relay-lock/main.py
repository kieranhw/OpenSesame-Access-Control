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

print("Starting OpenSesame Relay Lock...")
connect_wifi(settings["ssid"], settings["password"])
start_mdns(
    hostname=settings["hostname"],
    instance=settings["instance_name"],
    port=settings["port"]
)
server = HTTPServer(socketpool.SocketPool(wifi.radio), settings["port"])

lock_state = {"state": "unknown"}

def setup_output(pin_number):
    pin = digitalio.DigitalInOut(getattr(board, "GP{}".format(pin_number)))
    pin.direction = digitalio.Direction.OUTPUT
    pin.value = False # default to low
    return pin

red_led = setup_output(gpio["output"]["red_led"])
yellow_led = setup_output(gpio["output"]["yellow_led"])
green_led = setup_output(gpio["output"]["green_led"])
lock_relay = setup_output(gpio["output"]["lock_relay"])

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


def flash_leds(leds, flashes=10, delay=0.1):
    """
    Flash one or more LEDs to indicate failure.
    Args:
        leds (list): List of LED objects (e.g. [red_led, yellow_led]).
        flashes (int): Number of flashes (default 10).
        delay (float): Delay in seconds for on/off (default 0.1).
    """
    for _ in range(flashes):
        for led in leds:
            led.value = True
        time.sleep(delay)

        for led in leds:
            led.value = False
        time.sleep(delay)
        
    for led in leds:
        led.value = False

def request_handler(server, conn, method, path, body):
    if method == "POST" and path == "/lock":
        set_locked()
        return ("200 OK", "application/json", json.dumps(lock_state))

    elif method == "POST" and path == "/unlock":
        set_unlocked()
        return ("200 OK", "application/json", json.dumps(lock_state))

    elif method == "POST" and path == "/command_failed":
        server.write_response(conn, json.dumps(lock_state), content_type="application/json")
        flash_leds([red_led])
        return None

    elif method == "GET" and path == "/status":
        return ("200 OK", "application/json", json.dumps(lock_state))

    elif method == "GET" and path == "/info":
        mac_bytes = wifi.radio.mac_address
        mac_hex = str.capitalize(':'.join('{:02x}'.format(b) for b in mac_bytes))

        info = {
            "mac_address": mac_hex,
            "hostname": settings["hostname"],
            "instance_name": settings["instance_name"],
            "inputs": gpio["input"],
            "outputs": gpio["output"],
            "lock_state": lock_state
        }
        return ("200 OK", "application/json", json.dumps(info))
    else:
        return ("404 Not Found", "text/plain", "Not found")
    
try:
    server.serve_forever(request_handler)
    flash_leds([green_led, yellow_led])
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