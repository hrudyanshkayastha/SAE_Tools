# SAE Wazuh Build Evidence

## Build Command
```bash
wsl -u root bash -c "cd /home/nysro/sae_wazuh_build/src && make -j$(nproc) TARGET=server" 
```

## Execution Status
**State:** `FINISHED` 
**Exit Code:** `0` (Native Linux ext4 execution)

## Binary Verification
The following expected binaries physically exist on the native Linux filesystem:
- `/home/nysro/sae_wazuh_build/src/build/engine/wazuh-engine`
- `/home/nysro/sae_wazuh_build/src/build/bin/wazuh-manager-authd`
- `/home/nysro/sae_wazuh_build/src/build/bin/wazuh-manager-db`
- `/home/nysro/sae_wazuh_build/src/build/bin/wazuh-manager-remoted`

**Executable Test (`-h` output verification):**
All four binaries correctly execute natively and dump their expected help/startup configuration outputs.

## Version Verification
- Source tree `VERSION.json`: `v5.1.0 alpha0`
- Git commit tree: `8d0330c01` (Merge 5.0.1 into main)
- Binary stdout: `Wazuh v5.1.0 - Wazuh Inc.`
The inconsistency is resolved: the branch is 5.1.0 alpha.

## SAE End-to-End Ingestion Validation
A controlled test alert was pushed to the Wazuh alert log.
**Result from SAE Go Daemon:**
```text
Spawning Wazuh JSON listener... OK
SAE INGESTION SUCCESS: Normalized Wazuh Event -> Message: SAE Realtime E2E Test, Severity: Low, Product: SAE Tool (Wazuh Engine)
```
The real-time telemetry successfully traverses from the Wazuh log -> Go Native Consumer -> OCSF Mapper -> Backend console.
