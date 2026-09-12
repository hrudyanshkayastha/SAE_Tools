# Chunk E: UEBA v2 (Behavioral Analytics)

## 1. Goal
Advance the User and Entity Behavior Analytics engine (UEBA) to support multi-factor magnitude calculations (Risk Scoring) beyond simple binary anomaly triggers.

## 2. Inspection & Test
- Inspected `backend/internal/ueba/ueba.go`.
- Modified `EvaluateDeviation()` to calculate a standard deviation magnitude factor (`magnitude = (count - mean) / stddev`), clamping it into a 1-100 `RiskScore`.
- Injected `RiskScore` into the native OCSF payload.
- All UEBA tests natively compile and pass (`191/191` system tests).

## 3. Status
**VERIFIED**

## 4. Evidence
The native Go algorithm computes Risk Magnitude cleanly:
`[UEBA v2] ANOMALY DETECTED for User admin: 10 events this minute. Risk Score: 20 (Baseline Mean: 2.00, StdDev: 4.00)`
