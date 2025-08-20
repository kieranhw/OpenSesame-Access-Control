This directory contains the **OpenSesame Embedded Builder** that packages device source files into  `.zip` bundles for CircuitPython devices.

# Build Device Targets
-------------------

### 1\. Run the build script

From the project root, run:
```
cd os-device
python build.py
```
  
You’ll see a menu of available targets:

```
Step 1: Select a target to build
[1] opensesame_keypad
```
  
Enter the number of the target you want to build, then you will be prompted with the input steps.  

### 2\. Enter Wi-Fi credentials

The script will prompt you for:
*  **Wi-Fi SSID**
*  **Wi-Fi password**

These values will be appended to the target’s `settings.py` file automatically.

### 3\. Output

When the build completes, you’ll find a `.zip` file in the `/build` directory:
`/build/opensesame_keypad.zip`

Unzip this file and copy its contents to your CircuitPython device (e.g. Raspberry Pi Pico W).

# Creating New Device Programs
----------------------

To add a new device target:
1.  **Create a new source directory** under `src/access/` or `src/entries`, respectively. E.g. `src/access/opensesame-doorlock/`

2.  **Each device must include:**

* A `main.py` entrypoint file.
* A `settings.py` file with at least the following structure:
```
settings = {
	"hostname": "opensesame-keypad", # mDNS hostname
	"instance_name": "OpenSesame Keypad", # Friendly name
	"port": 80 # HTTP server port
}
```

* An optional `gpio.py` file which specifies pins that the user can override based on their wiring. This can contain development defaults, but will be overridden in the build output.

```
gpio = {
    "input": {
        "SENSOR": 1,
    },
    "output": {
        "LED": 2,
    }
}
```

Wi-Fi SSID and password will be appended by `build.py` during the build process.

3.  **Connect to the Wi-Fi and broadcast mDNS:**

At the start of the program, ensure to import and call both the settings and mDNS library in order to advertise the device to the OpenSesame hub.

```
from mdns import connect_wifi, start_mdns
from settings import settings

connect_wifi(settings["ssid"], settings["password"])
start_mdns(
	hostname=settings["hostname"],
	instance=settings["instance_name"],
	port=settings["port"]
)
```

4.  **Use the GPIO library to get the pins:**

To get the GPIO pins, in both development and production, you can use it in `main.py` like this:

```python
import board
import digitalio
from gpio import gpio

# Map GPIO numbers into variables
SENSOR_GPIO = gpio["input"]["SENSOR"]
LED_GPIO = gpio["output"]["LED"]

# Setup pins
sensor = digitalio.DigitalInOut(getattr(board, f"GP{SENSOR_GPIO}"))
sensor.direction = digitalio.Direction.INPUT
sensor.pull = digitalio.Pull.UP  # optional for buttons/sensors

led = digitalio.DigitalInOut(getattr(board, f"GP{LED_GPIO}"))
led.direction = digitalio.Direction.OUTPUT

# Example usage
if not sensor.value:
    led.value = True   # turn LED on
else:
    led.value = False  # turn LED off
```

5.  **Register the new target in `build.py`:**

```
TARGETS = {
	"opensesame_keypad": "src/access/opensesame-keypad",
	"opensesame_relay_lock": "src/access/opensesame-relay-lock", # new target
}
```

5.  **Build it** by running `python build.py` and selecting the new target.
