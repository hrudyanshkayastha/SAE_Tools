# SAE Falco Source Audit

## 1. Identity
* **Tool Name**: Falco (Cloud Native Runtime Security)
* **Version**: `0.45.0-rc1` (derived from git tags `0.45.0-rc1-16-g6b531e52`)
* **Repository**: Restored to pristine upstream state within `SAE/containers_security/`.
* **Commit**: `6b531e52cc030729e2ea4d1f184f1b553a5ad041` (Sep 10, 2026 - chore(cmake): bump libs to `0.26.0-rc2`)

## 2. Source Layout
The source resides strictly untouched within the SAE project root on the Windows NTFS volume:
`E:\New folder\SAE_Tools\SAE\containers_security\`

Previous destructive mass search-and-replace scripts have been completely reverted (`git reset --hard; git clean -fd`). No source files, CMake configurations, or internal identifiers have been modified. 
All compilations will occur strictly in an ephemeral WSL workspace to preserve this authoritative tree.

## 3. Build System
* **Build System**: CMake (C++)
* **Dependencies**: `cmake`, `make`, `gcc`/`g++`, `libssl-dev`, `libcurl4-openssl-dev`, `libyaml-cpp-dev`, `libjq-dev`, `libgrpc++-dev`, `protobuf-compiler`, `zlib1g-dev`, `libelf-dev` (for eBPF/kmod).
* **Runtime Subsystems**: Requires either the Falco kernel module (kmod), modern eBPF probe, or userspace instrumentation.

## 4. Integrity Assertion
The upstream Falco namespace has been 100% preserved. Falco will operate as a standalone binary outputting structured telemetry, avoiding any messy C++ linking or namespace collisions with the SAE Go backend.
