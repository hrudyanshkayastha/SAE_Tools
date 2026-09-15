# SAE TheHive License Audit

## 1. Upstream Identity
- **Repository**: TheHive (TheHive-Project)
- **Commit**: `d390a031c6a2e4e049969623e160a0a55e2dbd73`

## 2. License Details
- **Primary License**: GNU Affero General Public License v3.0 (AGPL-3.0)
- **Copyright Holders**: TheHive Project

## 3. SAE Distribution Posture
- **Modifications**: None. All corrupted internal naming modifications have been reversed to upstream AGPL-3.0 equivalents.
- **Integration Boundary**: TheHive functions as an external case management engine interacting with SAE via REST APIs or message queues. SAE does NOT statically or dynamically link to TheHive's Scala codebase, preserving strict process and memory separation. This architecture mitigates AGPL contamination risks for SAE's core. 
- **Legal Status**: Technical license audit only; not formal legal clearance.
