# SAE Trivy License Audit

## 1. Upstream Identity
- **Repository**: Trivy (Aqua Security)
- **Commit**: `ff327e75b70dab9e55f1fa4c5bb57a069ee5a06e`

## 2. License Details
- **Primary License**: Apache License 2.0
- **Copyright Holders**: Aqua Security Software Ltd.
- **Required Attribution**: Proper attribution to Aqua Security is required in redistributed materials.

## 3. SAE Distribution Posture
- **Modifications**: No upstream modifications are retained. All naming corruptions have been reversed.
- **Integration Boundary**: Trivy acts as a purely standalone executable. SAE calls the executable and reads its JSON output (via STDOUT or file). There is no static linking of Trivy libraries into SAE Core, thereby avoiding complex license entanglements.
- **Legal Status**: Technical license audit only; not formal legal clearance.
