# SAE Suricata Source Audit

## 1. Identity
* **Tool Name**: Suricata
* **Version**: 9.0.0-dev
* **Repository**: Found previously mapped in the workspace, now restored to pristine upstream state.
* **Commit**: `928ac012156fb8d393ce5ac4a496fde3c2e87b00` (Aug 6, 2026 - ssl: fix SSLv2 CLIENT_HELLO underflow)

## 2. Source Layout
The source resides untouched and pristine within the SAE project root on the Windows NTFS volume:
`E:\New folder\SAE_Tools\SAE\network_detection\`

No source files, `Makefile.am`, or `configure.ac` inside this authoritative tree were modified. The native build process (which required patching out Sphinx docs to bypass an automake crash) was conducted purely in an ephemeral WSL Linux ext4 workspace (`/home/nysro/sae_suricata_build/`) to ensure the official source tree remains cleanly verified.

### Key Internal Directories
- `/src/` - Primary C sources and packet processing logic
- `/rust/` - Rust-based parsers and app-layer handlers
- `/rules/` - Suricata rule definitions
- `/ebpf/` - eBPF bypass and acceleration
- `/python/` - Suricata-Update Python package

## 3. Build System
* **Primary Build System**: GNU Autotools (`autogen.sh`, `configure`, `make`)
* **Rust Build System**: Cargo (integrated into Autotools via `make`)
* **Dependencies**: `libpcre3-dev`, `libyaml-dev`, `libjansson-dev`, `libnspr4-dev`, `libnss3-dev`, `liblz4-dev`, `cargo`, `cbindgen`, `libcap-ng-dev`, `libnet1-dev`.

## 4. Integrity Assertion
Previous erroneous mass search-and-replace attempts that corrupted Suricata files, headers, and Makefiles have been successfully reverted via `git reset --hard` and `git clean -fd`. The upstream identifiers are now perfectly preserved in accordance with integration rules.
