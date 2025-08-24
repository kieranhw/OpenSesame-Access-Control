gpio = {
    "input": {
        # no inputs for this device
    },
    "output": {
        "red_led": 11,     # failed command indicator
        "yellow_led": 12,  # unlocked indicator
        "green_led": 13,   # locked indicator
        "lock_relay": 14  # single GPIO controlling the relay
    }
}