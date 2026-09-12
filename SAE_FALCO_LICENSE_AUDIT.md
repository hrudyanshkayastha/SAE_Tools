# SAE Falco License Audit

## 1. Upstream License
* **Software**: Falco
* **License**: Apache License, Version 2.0
* **Copyright**: The Falco Authors

## 2. Requirements & Restrictions
Under the Apache License 2.0:
* Commercial use, modification, distribution, and patent use are permitted.
* Modified files must include a notice of changes (if SAE were to distribute modified source directly).
* The `LICENSE` and `COPYING` files, as well as `NOTICE` files (if present), must be preserved and included in redistributed forms.
* Falco dependencies include libs (falcosecurity/libs) and various third-party components which are largely Apache 2.0 or MIT.

## 3. Compliance Measures Taken
* The `LICENSE` and `COPYING` files have been strictly preserved in `SAE/containers_security/`.
* No source identifiers have been stripped.
* SAE integrates Falco loosely via a decoupled JSON telemetry adapter (`eve.json` / `falco_events.json`), avoiding direct static/dynamic linking constraints, satisfying SAE architecture rules.

> **Disclaimer**: This audit identifies technical license boundaries. It does not constitute formal legal clearance.
