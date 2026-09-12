# Chunk A: Falco Runtime Validation

## 1. Goal
Achieve real Linux runtime validation for Falco kernel-level event capture.

## 2. Inspection & Test
- Environment check for Linux kernel headers and eBPF capabilities.
- The host OS remains Windows NT (Docker Desktop). Native Linux kernel syscalls cannot be intercepted.

## 3. Status
**BLOCKED**

## 4. Evidence
No fabricated evidence injected. Native execution is physically blocked by the host OS constraints.
