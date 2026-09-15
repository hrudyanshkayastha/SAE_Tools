# SAE Cortex License Audit

## 1. Upstream Identity
- **Repository**: Cortex (TheHive-Project)
- **Commit**: `061d49931752e5d956a94fb2193b2a2656360c6d`

## 2. License Details
- **Primary License**: GNU Affero General Public License v3.0 (AGPL-3.0)
- **Copyright Holders**: TheHive Project

## 3. SAE Distribution Posture
- **Modifications**: None. The pristine upstream source code is unmodified.
- **Integration Boundary**: Cortex functions as a standalone microservice engine that performs active observable analysis. SAE will interact with Cortex exclusively via its documented REST APIs (or Webhooks/JSON payloads). SAE does NOT statically or dynamically link to Cortex's Scala codebase, preserving strict process and memory separation to avoid AGPL contamination of SAE's core backend components.
- **Legal Status**: Technical license audit only; not formal legal clearance.
