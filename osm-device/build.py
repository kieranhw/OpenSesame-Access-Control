import os
import shutil
import zipfile
from pathlib import Path

PROJECT_FILES = [
    "common/mdns.py",
    "opensesame_keypad.py",
]
BUILD_DIR = Path("build")
OUTPUT_ZIP = Path("pico_app.zip")

def main():
    # Clean build dir
    if BUILD_DIR.exists():
        shutil.rmtree(BUILD_DIR)
    BUILD_DIR.mkdir(parents=True)

    # Ask user for Wi-Fi credentials
    ssid = input("Enter Wi-Fi SSID: ")
    password = input("Enter Wi-Fi password: ")

    # Write secrets.py
    secrets_content = f'''secrets = {{
    "ssid": "{ssid}",
    "password": "{password}",
}}
'''
    with open(BUILD_DIR / "secrets.py", "w") as f:
        f.write(secrets_content)

    # Copy project files
    for file in PROJECT_FILES:
        src = Path(file)
        dst = BUILD_DIR / src.name
        shutil.copy(src, dst)

    # Create zip archive
    with zipfile.ZipFile(OUTPUT_ZIP, "w") as zf:
        for file in BUILD_DIR.iterdir():
            zf.write(file, file.name)

    print(f"✅ Build complete: {OUTPUT_ZIP}")
    print("👉 To deploy: unzip and copy contents to CIRCUITPY drive.")


if __name__ == "__main__":
    main()