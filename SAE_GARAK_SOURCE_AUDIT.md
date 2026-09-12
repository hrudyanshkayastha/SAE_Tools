# SAE Garak Source Audit

## 1. Upstream Identity
- **Project**: NVIDIA Garak (NVIDIA/garak)
- **Engine**: Python native library installed on the host OS.
- **Version**: 0.16.0 (via `pip list`)

## 2. Integration Type
As directed by the requirement to adhere strictly to upstream sources and avoiding unneeded backends, Garak is integrated using the globally available Python environment natively accessible by the SAE Go backend. 
- The native command line utility `garak` was utilized.
- No Garak core source code was modified.
