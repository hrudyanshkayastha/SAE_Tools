import json
import subprocess
import time
import urllib.request
import urllib.parse
files = ['attack.json', 'benign.json', 'benign.json', 'benign.json', 'zeek.json', 'suri.json']
for f in files:
    with open('/mnt/e/New folder/SAE_Tools/SAE/' + f, 'r') as file:
        data = file.read()
        subprocess.run(['redis-cli', 'XADD', 'sae_events', '*', 'data', data])
        time.sleep(1)
