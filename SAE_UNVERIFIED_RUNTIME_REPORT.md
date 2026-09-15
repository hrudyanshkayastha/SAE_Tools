# SAE Blocked Runtimes Report

In accordance with strict verification rules, no synthetic runtime evidence has been fabricated. The following components are structurally integrated and pass all unit tests but remain **BLOCKED** from physical E2E verification due to precise environment limitations.

## 1. Falco (SAE Containers Security)
1. **Upstream Identity**: falcosecurity/falco
2. **Integrated Commit**: 6b531e52cc030729e2ea4d1f184f1b553a5ad041
3. **SAE Role**: Container/Host runtime intrusion detection.
4. **Runtime Dependencies**: Native Linux kernel (CONFIG_FTRACE enabled for eBPF or loadable kernel modules).
5. **Exact Blocker**: The local environment runs on a Windows Host OS with a WSL2 kernel (6.6.87.2-microsoft-standard-WSL2) that lacks trace/tracepoints, making syscall interception physically impossible.
6. **Final Status**: ?? **BLOCKED**

## 2. KubeArmor (SAE CloudGuard)
1. **Upstream Identity**: kubearmor/KubeArmor
2. **Integrated Commit**: cdb06520561dc2db4cb25d9e1c4f555ac733de1
3. **SAE Role**: Kubernetes policy enforcement (LSM/eBPF).
4. **Runtime Dependencies**: A live Kubernetes cluster running on a Linux node with AppArmor, SELinux, or BPF-LSM enabled.
5. **Exact Blocker**: No local Kubernetes cluster is active (Docker Desktop is offline/unavailable in the current environment context), and WSL2 cannot easily mock K8s LSM bindings without a dedicated host.
6. **Final Status**: ?? **BLOCKED**

## 3. ScoutSuite (SAE Cloud Assessment)
1. **Upstream Identity**: nccgroup/ScoutSuite
2. **Integrated Commit**: 7909f2fc6186063e5c9e7ddef8c4d7d1072c8f3d
3. **SAE Role**: Cloud configuration assessment and posture management.
4. **Runtime Dependencies**: Authorized, read-only cloud credentials (AWS IAM, GCP Service Account, Azure AD).
5. **Exact Blocker**: The current execution environment does not possess valid cloud credentials. Fabricating a cloud assessment payload or mocking the API violates SAE verification rules.
6. **Final Status**: ?? **BLOCKED**

---

### SAE Integration Confidence
- **Adapters + OCSF Mapping**: ? Implemented
- **SAE Test Suite**: ? 195/195 passing
- **Fake/Synthetic Evidence**: ? Not Used
