# SAE Suricata Build Evidence

## 1. Source Identity
- **Tool**: Suricata
- **Version**: `9.0.0-dev`
- **Commit**: `928ac012156fb8d393ce5ac4a496fde3c2e87b00`
- **Branch**: `master`

## 2. Build Environment
- **OS**: Ubuntu 26.04 via WSL2 (native ext4 filesystem)
- **Build Location**: `/home/nysro/sae_suricata_build/`
- **Compiler**: GCC 15.2.0, Rustc 1.93.1, Cargo 1.93.1

## 3. Prerequisites Installed (APT)
- `libpcre2-dev`
- `libyaml-dev`
- `libjansson-dev`
- `libnspr4-dev`
- `libnss3-dev`
- `liblz4-dev`
- `cargo` & `rustc`
- `cbindgen`
- `libcap-ng-dev`
- `libnet1-dev`
- `autoconf` & `automake` & `libtool`

## 4. Compilation Commands
```bash
dos2unix ./*
./autogen.sh
./configure --disable-shared
make -j4
```

## 5. Build Status
- **Exit Code**: `0`
- **Output Binary Path**: `/home/nysro/sae_suricata_build/src/suricata`
- **Binary Size**: `117M`
- **Version Validated via Binary**:
  ```text
  This is Suricata version 9.0.0-dev (928ac0121 2026-09-05)
  Features: PCAP_SET_BUFF AF_PACKET HAVE_PACKET_FANOUT LIBCAP_NG LIBNET1.1 HAVE_HTP_URI_NORMALIZE_HOOK PCRE_JIT HAVE_NSS HTTP2_DECOMPRESSION HAVE_LUA HAVE_JA3 HAVE_JA4 HAVE_LIBJANSSON UNIX_SOCKET TLS TLS_C11 RUST POPCNT64 
  ```

## 6. Pre-Build Gate Validations
All tests passed, verifying that Suricata's Go adapter successfully maps `eve.json` events without breaking Wazuh or Zeek.
```
=== RUN   TestMapEVEToOCSF_Alert — PASS
=== RUN   TestMapEVEToOCSF_Flow — PASS
... (All 7 Suricata tests passing)
```
**Total Backend Tests Passing: 15/15** (Wazuh: 4, Zeek: 4, Suricata: 7).
