# SAE License Registry & Attributions

## Proprietary Code
The `SAE Core` (Written natively in Go and Python, mapped within `sae-core/` and the future agent/detection folders) is strictly proprietary software owned by SAE. 

## Third-Party Open Source Components

In accordance with open-source legal requirements, the following components are utilized within the SAE product ecosystem. They are executed as isolated network services and are not statically linked to the SAE proprietary codebase.

### 1. Wazuh Engine
- **Original Source:** Wazuh Inc. (Fork of OSSEC)
- **License:** GNU General Public License version 2.0 (GPLv2)
- **Location:** `SAE/engine/` and `SAE/licenses/WAZUH_GPLv2.txt`
- **Integration Boundary:** The engine is compiled as a standalone application. SAE interacts with it strictly via REST API (Port 55000) and by tailing its filesystem JSON outputs. No GPLv2 C/C++ code is compiled into the SAE Go binaries.

### 2. OCSF (Open Cybersecurity Schema Framework)
- **Original Source:** OCSF Community
- **License:** Apache 2.0
- **Location:** `sae-core/models/ocsf.go`
- **Integration Boundary:** Struct schemas were utilized to build the unified event model. 
