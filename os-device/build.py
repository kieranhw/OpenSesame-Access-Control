import shutil
import zipfile
from pathlib import Path
import os
import sys

# cd to build.py
SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
os.chdir(SCRIPT_DIR)
sys.path.insert(0, SCRIPT_DIR)

BUILD_ROOT = Path("build")
BOLD = "\033[1m"
GREEN = "\033[32m"
RESET = "\033[0m"

# specify targets for each device build output, where the key is
# the name of the output zip and the value is the src directory
TARGETS = {
    "opensesame_keypad": "src/access/opensesame-keypad",
    "opensesame_relay_lock": "src/entry/opensesame-relay-lock",
}
COMMON_PATH = "src/common"

def append_wifi_settings(settings_file: Path, ssid: str, password: str):
    with settings_file.open("a", encoding="utf-8") as f:
        f.write(
            f'\n# Wi-Fi settings\n'
            f'settings["ssid"] = "{ssid}"\n'
            f'settings["password"] = "{password}"\n'
        )

def copy_all_files(src_dir: Path, dest_dir: Path):
    if not src_dir.exists():
        print(f"Skipping missing directory: {src_dir}")
        return

    for root, _, files in os.walk(src_dir):
        for file in files:
            src_file = Path(root) / file
            rel_path = src_file.relative_to(src_dir)
            dest_file = dest_dir / rel_path
            dest_file.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(src_file, dest_file)


def parse_gpio_file(gpio_file: Path) -> dict:
    gpio = {}
    if not gpio_file.exists():
        return gpio

    namespace = {}
    with open(gpio_file, "r", encoding="utf-8") as f:
        code = f.read()
        exec(code, namespace)

    if "gpio" in namespace and isinstance(namespace["gpio"], dict):
        gpio = namespace["gpio"]

    return gpio

def override_gpio_pins(gpio_file: Path, overrides: dict):
    gpio = parse_gpio_file(gpio_file)
    if not gpio:
        return

    # Apply overrides
    for section, pins in overrides.items():
        if section not in gpio:
            gpio[section] = {}
        for key, new_val in pins.items():
            gpio[section][key] = new_val

    # Rewrite gpio.py
    with open(gpio_file, "w", encoding="utf-8") as f:
        f.write("# Auto-generated with user overrides\n\ngpio = {\n")
        for section, pins in gpio.items():
            f.write(f'    "{section}": {{\n')
            for key, val in pins.items():
                f.write(f'        "{key}": {val},\n')
            f.write("    },\n")
        f.write("}\n")

def build_for_target(
    target_name: str,
    src_dir: str,
    common_dir: str,
    ssid: str,
    password: str,
    gpio_overrides: dict,
):
    build_dir = BUILD_ROOT / target_name
    if build_dir.exists():
        shutil.rmtree(build_dir)
    build_dir.mkdir(parents=True)

    # copy target directory files
    copy_all_files(Path(src_dir), build_dir)

    # copy common files
    if common_dir:
        copy_all_files(Path(common_dir), build_dir)

    # ensure settings.py exists
    settings_file = build_dir / "os_settings.py"
    if not settings_file.exists():
        raise FileNotFoundError(
            f"Target {target_name} is missing a os_settings.py file in {src_dir}"
        )

    # append Wi-Fi settings
    append_wifi_settings(settings_file, ssid, password)

    # apply GPIO overrides
    gpio_file = build_dir / "os_gpio.py"
    if gpio_file.exists():
        override_gpio_pins(gpio_file, gpio_overrides)

    # create a zip containing each dependency in the target
    target_zip = BUILD_ROOT / f"{target_name}.zip"
    if target_zip.exists():
        target_zip.unlink()

    with zipfile.ZipFile(target_zip, "w") as zipf:
        for root, _, files in os.walk(build_dir):
            for file in files:
                file_path = Path(root) / file
                rel_path = file_path.relative_to(build_dir)
                zipf.write(file_path, rel_path)

    print(f"-> Built {target_name} to {target_zip}")

def main():
    BUILD_ROOT.mkdir(parents=True, exist_ok=True)
    print(f'''
 _____           _              _     _          _ 
| ____|_ __ ___ | |__   ___  __| | __| | ___  __| |
|  _| | '_ ` _ \| '_ \ / _ \/ _` |/ _` |/ _ \/ _` |
| |___| | | | | | |_) |  __/ (_| | (_| |  __/ (_| |
|_____|_| |_|_|_|_.__/ \___|\__,_|\__,_|\___|\__,_|
| __ ) _   _(_) | __| | ___ _ __                   
|  _ \| | | | | |/ _` |/ _ \ '__|                  
| |_) | |_| | | | (_| |  __/ |                     
|____/ \__,_|_|_|\__,_|\___|_|                                                                                                                                                                                                                             
          ''')

    print(f"\n{BOLD}OpenSesame Embedded Builder{RESET}\n")
    print(f"{BOLD}Step 1: Select a target to build{RESET}\n")
    targets = list(TARGETS.keys())
    for i, t in enumerate(targets, start=1):
        print(f"[{i}] {t}")
    choice = input("\nEnter the number of the target: ")

    try:
        choice_idx = int(choice) - 1
        target_name = targets[choice_idx]
    except (ValueError, IndexError):
        print("Invalid choice. Exiting.")
        sys.exit(1)

    src_dir = TARGETS[target_name]

    # get Wi-Fi credentials
    print(f"\n{BOLD}Step 2: Enter Wi-Fi settings{RESET}\n")
    ssid = input("Enter Wi-Fi SSID: ")
    password = input("Enter Wi-Fi password: ")

    # GPIO overrides
    print(f"\n{BOLD}Step 3: Configure GPIO pins{RESET}")
    gpio_file = Path(src_dir) / "os_gpio.py"
    gpio_defaults = parse_gpio_file(gpio_file)
    gpio_overrides = {"input": {}, "output": {}}

    if gpio_defaults:
        for section in ["input", "output"]:
            if section in gpio_defaults:
                if not gpio_defaults[section]:
                    # no gpio for this section, skip
                    continue
                
                print(f"\nConfigure {section.upper()} pins:")
                for key, default_val in gpio_defaults[section].items():
                    new_val = input(f"{key} (default {default_val}): ").strip()
                    if new_val:
                        try:
                            gpio_overrides[section][key] = int(new_val)
                        except ValueError:
                            print(f"Invalid input for {key}, exiting")
                            sys.exit(1)

    print("\nBuilding target...")
    build_for_target(target_name, src_dir, COMMON_PATH, ssid, password, gpio_overrides)

    print(
        f"\n{BOLD}{GREEN}Build complete! You may now load the module from /build to your device.{RESET}\n"
    )

if __name__ == "__main__":
    main()