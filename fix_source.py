import os
import re

src_dir = r"E:\New folder\SAE_Tools\SAE\engine"

replacements = {
    "SAE TOOL": "WAZUH",
    "Sae Tool": "Wazuh",
    "SAE Tool": "Wazuh",
    "sae tool": "wazuh",
    "sae_tool": "wazuh",
    "SAE_TOOL": "WAZUH"
}

def is_text(filename):
    try:
        with open(filename, 'tr', encoding='utf-8') as check_file:
            check_file.read(1024)
            return True
    except:
        return False

print("Starting restoration of C/C++ variables...")
for root, dirs, files in os.walk(src_dir):
    if '.git' in dirs:
        dirs.remove('.git')
    for name in files:
        filepath = os.path.join(root, name)
        if not is_text(filepath):
            continue
        try:
            with open(filepath, 'r', encoding='utf-8') as f:
                content = f.read()
            
            new_content = content
            for old, new in replacements.items():
                new_content = new_content.replace(old, new)
            
            if new_content != content:
                with open(filepath, 'w', encoding='utf-8') as f:
                    f.write(new_content)
        except Exception as e:
            pass

print("Done replacing text.")
