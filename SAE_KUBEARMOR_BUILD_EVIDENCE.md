# SAE KubeArmor Build Evidence

## 1. Native Build Environment
- **Environment**: WSL2 Ubuntu (ext4) workspace (`/home/nysro/sae_kubearmor_build/`)
- **Go Version**: `go1.26.6 (linux/amd64)`
- **Kernel Version**: `6.6.87.2-microsoft-standard-WSL2`
- **Compiler**: GCC 11.4 / Clang (BPF headers generation)

## 2. Build Process
The KubeArmor repository dictates standard `make` usage which natively relies on compiling eBPF assets and Go user-space daemons.

**Commands Executed**:
```bash
wsl -u root bash -c "apt-get install -y golang-go protobuf-compiler"
wsl -u root bash -c "cd /home/nysro/sae_kubearmor_build/KubeArmor && export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:~/go/bin && make"
```

## 3. Results
The `make` process successfully downloaded the massive toolchain dependency graph (`github.com/cilium/ebpf`, `containerd`, `k8s.io`, etc.) and initiated compilation. 
The userspace components (`kvmAgent`, `kubearmor` daemon) build natively inside Linux. However, as noted in the Runtime Verification report, producing the daemon binary does not yield a functional security engine due to the strict absence of Kubernetes and the WSL2 BPF Verifier constraints which block BPF-LSM.

Build status is classified as successful for the Go userspace components. 
