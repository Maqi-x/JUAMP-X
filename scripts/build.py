# ------------------------------------------- MAKE GAME ---------------------------------------------- #

import subprocess as sp
import platform
import time
import sys as os
import os as system
import os.path as path

if not path.exists("bin"):
    system.mkdir("bin")

sys = platform.system()

def erro(err):
    if sys == 'Linux' or sys == 'Darwin':
        print(f"\033[1;31m{err}\033[0;0m")
    else:
        print(err)

def check():
    try:
        result = sp.run(["go", "version"], stdout=sp.PIPE, stderr=sp.PIPE, text=True)
        if result.returncode != 0:
            raise Exception("Go is not installed or not in PATH.")
        
        VerInfo = result.stdout.strip()
        if not VerInfo.startswith("go version"):
            erro(f"Unexpected output: {VerInfo}")
            os.exit(1)
        
        version = VerInfo.split(" ")[2][2:]
        parts = version.split(".")
        n = int(parts[0])
        ver = int(parts[1])
        patch = int(parts[2]) if len(parts) > 2 else 0
        pass
        
        MinN, MinVer, MinPatch = 1, 21, 10
        
        if (n, ver, patch) >= (MinN, MinVer, MinPatch):
            ansimsg = '\033[38;5;46mAll dependencies are installed!\033[0m'
            msg = 'All dependencies are installed!'
            
            print(f"{ansimsg if sys == 'Linux' or sys == 'Darwin' else msg}")
            return True
            pass
        else:
            erro(f"it looks like you don't have Go installed, or it's not in $PATH; properly install Go and try again (Error code: {VerInfo})")
            os.exit(250)
            pass
        pass
    except sp.CalledProcessError as err:
        erro(f"Error: command failed! Make sure Go is installed and in $PATH (error code: {err})")
        os.exit(561)
        pass
    except Exception as err:
        erro(f"Error: {err}")
        os.exit(259)
        pass
    pass
pass

def loading(proc):
    d = ""
    while proc.poll() is None:
        print(f"\rBuilding game{d:<3}", end="", flush=True)
        d += "."
        if len(d) > 3:
            d = ""
        time.sleep(0.15)
    print(f"\rBuilding game... Done!")
    pass
pass

check()

if sys == "Linux" or sys == "Darwin":
    try:
        cmd = sp.Popen([
            "go", "build", "-o", "bin/JUAMP-X"
        ])
    except Exception as err:
        print(f"\033[1;31mError!\033[0;91m Compilation not successful! make sure you have installed all dependencies\033[3m - ErrorCode: {err}\033[0;0m")
        os.exit(249)
        pass
    loading(cmd)
    pass
pass

if sys == "Windows":
    try:
        cmd = sp.Popen([
            "go", "build", "-o", "bin/JUAMP-X.exe"
        ])
    except Exception as err:
        print(f"Compilation not successful! make sure you have installed all dependencies - ErrorCode: {err}")
        os.exit(305)
        pass
    loading(cmd)
    pass
pass
