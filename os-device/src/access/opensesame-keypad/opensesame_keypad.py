import machine
import time
from matrix_keypad import MatrixKeypad

if __name__ == "__main__":
    # Define row pins (adjust these to match your wiring)
    row_pins = [
        machine.Pin(9, machine.Pin.IN, machine.Pin.PULL_UP),   # Row 0: "1", "2", "3"
        machine.Pin(4, machine.Pin.IN, machine.Pin.PULL_UP),   # Row 1: "4", "5", "6"
        machine.Pin(5, machine.Pin.IN, machine.Pin.PULL_UP),   # Row 2: "7", "8", "9"
        machine.Pin(7, machine.Pin.IN, machine.Pin.PULL_UP)    # Row 3: "*", "0", "#"
    ]

    # Define column pins (adjust these to match your wiring)
    col_pins = [
        machine.Pin(8, machine.Pin.OUT),    # Column 0
        machine.Pin(10, machine.Pin.OUT),   # Column 1
        machine.Pin(6, machine.Pin.OUT)     # Column 2
    ]

    # Define the keypad layout:
    keys = (
        ("1", "2", "3"),
        ("4", "5", "6"),
        ("7", "8", "9"),
        ("*", "0", "#")
    )

    keypad = MatrixKeypad(row_pins, col_pins, keys, debounce_ms=50)

    while True:
        new_keys = keypad.rising_edges
        if new_keys:
            print("New key pressed:", new_keys)
        time.sleep(0.1)