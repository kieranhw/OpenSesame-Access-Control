# Default GPIO pin mappings
gpio = {
    "input": {
        "DOOR_SENSOR": 14,   # reed switch or magnetic sensor
        "BUTTON": 11,        # push button
    },
    "output": {
        "LOCK": 15,          # relay to lock door
        "UNLOCK": 12,        # relay to unlock door
        "LED": 2,            # status LED
    }
}