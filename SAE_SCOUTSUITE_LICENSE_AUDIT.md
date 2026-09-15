# SAE ScoutSuite License Audit

## 1. Upstream Identity
- **Repository**: ScoutSuite (NCC Group)
- **Commit**: `7909f2fc6186063e5c9e7ddef8c4d7d1072c8f3d`

## 2. License Details
- **Primary License**: GPL v2.0 (General Public License version 2.0)
- **Copyright Holders**: NCC Group
- **Attribution**: Must retain the original NCC Group attributions and license texts.

## 3. SAE Distribution Posture
- **Modifications**: None. All corrupted internal naming modifications have been reversed.
- **Integration Boundary**: Because ScoutSuite is a GPLv2 tool (written in Python) and SAE Core is a proprietary/internal engine (written in Go), integrating them directly via shared memory or library calls would potentially taint the SAE Core under the GPL's viral clause. To avoid this, SAE enforces absolute process isolation. ScoutSuite is executed purely as a standalone executable CLI tool, communicating with SAE solely via the standardized JSON reports it writes to disk.
- **Legal Status**: Technical license audit only; not formal legal clearance.
