# SAE Project Structure

## Overview
The SAE (Security AI Engine) has been structurally unified into a single root directory: `E:\New folder\SAE_Tools\SAE\`. 
To preserve the complex make-based compilation paths of the integrated Wazuh engine while satisfying the logical SAE component structure, we utilized NTFS Junctions to map the core capabilities securely.

## Logical Directory Tree

```text
E:\New folder\SAE_Tools\SAE\
├── core\            (Junction -> engine/src: The core C/C++ event engines)
├── api\             (Junction -> engine/api: The Python/Node REST API)
├── rules\           (Junction -> engine/ruleset/rules: Detection signatures)
├── decoders\        (Junction -> engine/ruleset/decoders: Log parsers)
├── configuration\   (Junction -> engine/etc: System templates and configs)
├── tests\           (Junction -> engine/tests: Integration and unit tests)
├── licenses\        (Contains GPLv2 and MIT licenses for third-party components)
├── agent\           (Reserved for native SAE agent Go code)
├── detection\       (Reserved for SAE OCSF correlation rules)
├── storage\         (Reserved for ClickHouse/PostgreSQL bindings)
├── frontend\        (Reserved for SAE UI)
├── integrations\    (Reserved for SAE orchestration webhooks)
└── engine\          (Physical isolated Wazuh compilation codebase)
```

## Architectural Rationale
By moving the Wazuh source code entirely into `SAE/engine` and mapping it logically, we achieved three things:
1. **One Root:** There is no longer a separate "Wazuh tool" floating alongside SAE.
2. **Build Safety:** The strict relative paths defined in the C/C++ Makefiles (`../etc`, `../ruleset`) are preserved physically.
3. **Abstraction:** The SAE Go backend can reference `SAE/api` and `SAE/rules` natively, treating the engine as a modular component rather than a standalone product.
