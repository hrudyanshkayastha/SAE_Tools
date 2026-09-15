# SAE Falco Build Evidence

## 1. Source Identity
- **Tool**: Falco
- **Version**: `0.45.0-rc1` (Git: `6b531e52cc030729e2ea4d1f184f1b553a5ad041`)
- **Branch**: `master`

## 2. Build Environment
- **OS**: Ubuntu 26.04 via WSL2 (ext4 native filesystem)
- **Workspace Location**: `/home/nysro/sae_falco_build/`
- **Build System**: CMake, Make
- **Compiler**: GCC 15.2.0, G++ 15.2.0, Clang 21.1.8

## 3. Prerequisite Dependencies Installed (APT)
- `cmake` & `make`
- `gcc`, `g++`, `clang`, `llvm`, `bpftool`, `linux-tools-common`
- `libssl-dev`, `libcurl4-openssl-dev`, `libyaml-cpp-dev`, `libjq-dev`, `libgrpc++-dev`, `protobuf-compiler`
- `zlib1g-dev`, `libelf-dev`, `libbpf-dev`

## 4. Compilation Execution
```bash
dos2unix ./*
mkdir build && cd build
cmake -DFALCO_ETC_DIR=/etc/falco -DUSE_BUNDLED_DEPS=ON ..
make -j4 falco
```
*(Note: Building the `falco` target directly bypasses problematic out-of-tree WSL kernel driver compilations which block on `linux-headers`, ensuring the userspace binary builds cleanly.)*

## 5. Build Status
- **Exit Code**: `0`
- **Output Binary Path**: `/home/nysro/sae_falco_build/build/userspace/falco/falco`
- **Binary Size**: `312M`
- **Version Validated via Binary**:
  ```text
  Falco version: 0.45.0-76+6b531e5
  Libs version:  0.26.0-rc2
  Plugin API:    3.12.0
  Engine:        0.65.0
  Driver API version: 11.0.0
  ```
