# SAE KubeArmor License Audit

## 1. Upstream Identity
- **Repository**: KubeArmor
- **Commit**: `fcdb06520561dc2db4cb25d9e1c4f555ac733de1`

## 2. License Details
- **Primary License**: Apache License 2.0
- **Copyright Holders**: KubeArmor Authors / AccuKnox
- **Required Attribution**: Standard Apache 2.0 attribution notices required in redistributed binaries/source.

## 3. SAE Distribution Posture
- **Modifications**: SAE has reverted all unauthorized global string replacements. The source is identical to upstream.
- **Integration**: SAE implements a loose coupling via `sae-core/internal/kubearmor/` adapter logic that ingests telemetry output, completely avoiding GPL or proprietary viral contamination between systems.
- **Legal Clearance**: Technical license audit only; not formal legal clearance.
