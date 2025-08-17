import shutil
import zipfile
from pathlib import Path

BOLD = "\033[1m"
GREEN = "\033[32m"
RESET = "\033[0m"

# specify targets for each device build output, where key is 
# the name of the output zip and values are the src dependency files
TARGETS = {
    "opensesame_keypad": [
        "src/access/opensesame-keypad/opensesame_keypad.py",
        "src/access/opensesame-keypad/matrix_keypad.py",
        "src/access/opensesame-keypad/mdns.py",
    ],
}

BUILD_ROOT = Path("build")
OUTPUT_ROOT = Path("build")
SECRETS = None

def make_secrets_content(ssid: str, password: str) -> str:
    return f'''secrets = {{
    "ssid": "{ssid}",
    "password": "{password}",
}}
'''

def build_for_target(target_name: str, dependencies: list[str], shared_ssid: str, shared_password: str):
    build_dir = BUILD_ROOT / target_name
    if build_dir.exists():
        shutil.rmtree(build_dir)
    build_dir.mkdir(parents=True)

    # create secrets file
    secrets_content = make_secrets_content(shared_ssid, shared_password)
    (build_dir / "secrets.py").write_text(secrets_content, encoding="utf-8")

    # copy required source files into per-target build dir
    for file in dependencies:
        src = Path(file)
        dst = build_dir / src.name
        if not src.exists():
            raise FileNotFoundError(f"Source file not found: {src}")
        shutil.copy(src, dst)

    # create a zip containing each dependency in the target
    target_zip = OUTPUT_ROOT / f"{target_name}.zip"
    if target_zip.exists():
        target_zip.unlink()

    with zipfile.ZipFile(target_zip, "w") as zip:
        for file in build_dir.iterdir():
            zip.write(file, file.name)

    # delete the non-zipped directory
    shutil.rmtree(build_dir)

    print(f"-> Built {target_name} to {target_zip}")


def main():
    OUTPUT_ROOT.mkdir(parents=True, exist_ok=True)

    print(f"\n{BOLD}OpenSesame Embedded Builder{RESET}")
    print(f"{BOLD}To build your files please first provide some required information.{RESET}\n")

    # shared Wi-Fi credentials for all targets
    print(f"{BOLD}Step 1: To enable mDNS, please provide Wi-Fi credentials for wireless devices.{RESET}\n")
    ssid = input("Enter Wi-Fi SSID: ")
    password = input("Enter Wi-Fi password: ")

    # build all targets
    print("\nBuilding targets...")
    for target, dependencies in TARGETS.items():
        build_for_target(target, dependencies, ssid, password)
        
    print(f"\n{BOLD}{GREEN}Build complete, you may now load the zips from /build to your device.{RESET}\n")
    
if __name__ == "__main__":
    main()