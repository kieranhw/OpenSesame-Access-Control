import machine
import time

class MatrixKeypad:
    def __init__(self, row_pins, col_pins, keys, debounce_ms=50):
        """
        row_pins: list of machine.Pin objects for rows (inputs with pull-up)
        col_pins: list of machine.Pin objects for columns (outputs)
        keys: tuple of tuples representing the keypad layout (rows x columns)
        debounce_ms: debounce time in milliseconds
        """
        self.rows = row_pins
        self.cols = col_pins
        self.keys = keys
        self.debounce_ms = debounce_ms

        # Initialise row pins as inputs with pull-ups (default high)
        for row in self.rows:
            row.init(mode=machine.Pin.IN, pull=machine.Pin.PULL_UP)

        # Initialise column pins as outputs and set them high (default state)
        for col in self.cols:
            col.init(mode=machine.Pin.OUT)
            col.value(1)

        # Store previously pressed keys for rising-edge detection.
        self._prev_keys = set()

    def scan(self):
        pressed = set()
        # For each column, set it low and check the rows.
        for col_index, col in enumerate(self.cols):
            col.value(0)
            time.sleep_ms(1)  # allow signal to settle

            # Check each row for a stable press.
            for row_index, row in enumerate(self.rows):
                # In active low mode, a pressed key will bring the row to 0.
                if row.value() == 0:
                    start = time.ticks_ms()
                    # Debounce: ensure the row remains low for debounce_ms.
                    while time.ticks_diff(time.ticks_ms(), start) < self.debounce_ms:
                        if row.value() != 0:
                            break
                        time.sleep_ms(1)
                    else:
                        pressed.add(self.keys[row_index][col_index])
                        # Only register one key per column.
                        break
            # Restore the column to high before moving to the next.
            col.value(1)
        return pressed

    @property
    def rising_edges(self):
        """
        Returns only keys that are newly pressed (rising edge).
        """
        current = self.scan()
        new_keys = current - self._prev_keys
        self._prev_keys = current
        return new_keys

# Example usage:
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