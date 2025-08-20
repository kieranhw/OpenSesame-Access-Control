# main.py
import time
import board
import digitalio
from settings import settings
from gpio import gpio
from mdns import connect_wifi, start_mdns

# Connect to Wi-Fi
connect_wifi(settings["ssid"], settings["password"])

# Start mDNS
start_mdns(
    hostname=settings["hostname"],
    instance=settings["instance_name"],
    port=settings["port"]
)

# Map GPIO numbers into variables
LED1_GPIO = gpio["output"]["LED1"]
LED2_GPIO = gpio["output"]["LED2"]
LED3_GPIO = gpio["output"]["LED3"]

# Setup pins
led1 = digitalio.DigitalInOut(getattr(board, f"GP{LED1_GPIO}"))
led1.direction = digitalio.Direction.OUTPUT

led2 = digitalio.DigitalInOut(getattr(board, f"GP{LED2_GPIO}"))
led2.direction = digitalio.Direction.OUTPUT

led3 = digitalio.DigitalInOut(getattr(board, f"GP{LED3_GPIO}"))
led3.direction = digitalio.Direction.OUTPUT

# Flash LEDs in sequence
while True:
    for led in [led1, led2, led3]:
        led.value = True
        time.sleep(0.3)
        led.value = False
        time.sleep(0.3)