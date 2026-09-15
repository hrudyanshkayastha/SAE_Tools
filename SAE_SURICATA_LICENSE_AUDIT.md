# SAE Suricata License Audit

## 1. Upstream License
* **Software**: Suricata Network IDS/IPS/NSM
* **License**: GNU General Public License v2 (GPLv2)
* **Copyright**: Open Information Security Foundation (OISF) and contributors.

## 2. Requirements & Restrictions
Under the GPLv2:
* Any modifications to the source code, if distributed, must also be open-sourced under the GPLv2.
* Attribution and copyright notices must remain intact.
* As SAE integrates Suricata via an external telemetry bridge (parsing `eve.json` natively) and does not directly link against Suricata's internal C/Rust libraries, this loosely coupling model satisfies SAE's architectural constraints without forcing the Go orchestration layer into a derivative GPL standing, pending legal review.

## 3. Compliance Measures Taken
* The `LICENSE` and `COPYING` files have been preserved unmodified.
* No source identifiers have been stripped.
* The `SAE/network_detection/` boundary acts as the strict separation line between OISF GPLv2 code and SAE proprietary wrapping.

> **Disclaimer**: This audit identifies technical license boundaries. It does not constitute formal legal clearance.
