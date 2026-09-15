# SAE Garak Test Report

## 1. Adapter Test Verification
The Go backend test suite (`backend/internal/garak/`) processed actual Garak JSONL structures.

| Target | Result | Coverage |
|---|---|---|
| `TestParseGarakEvalToOCSF_Empty` | PASS | Handles empty lines safely. |
| `TestParseGarakEvalToOCSF_InvalidJSON` | PASS | Safely catches broken JSON lines. |
| `TestParseGarakEvalToOCSF_NotEval` | PASS | Ignores setup/initialization JSON lines in the report. |
| `TestParseGarakEvalToOCSF_ValidEval` | PASS | Successfully parses a full `eval` JSON result block. |
| `TestParseGarakEvalToOCSF_SeverityCritical` | PASS | Maps DEFCON 1-2 to Critical OCSF severity. |
| `TestParseGarakEvalToOCSF_SeverityHigh` | PASS | Maps DEFCON 3 to High OCSF severity. |
| `TestParseGarakEvalToOCSF_SeverityMedium` | PASS | Maps DEFCON 4 to Medium OCSF severity. |
| `TestParseGarakEvalToOCSF_SeverityLow` | PASS | Maps DEFCON 5 (if failures exist) to Low OCSF severity. |
| `TestParseGarakEvalToOCSF_Observables` | PASS | Extracts detector names and plugin groups to observables. |
| `TestParseGarakEvalToOCSF_IgnoresSummary` | PASS | Safely bypasses `_summary` keys in the JSON tree. |
| `TestParseGarakEvalToOCSF_IgnoresProbes` | PASS | Safely bypasses internal probe keys in the JSON tree. |
| `TestParseGarakEvalToOCSF_MessageFormat` | PASS | Assembles descriptive OCSF Activity Messages containing pass/fail metrics. |
| `TestParseGarakEvalToOCSF_MultipleDetectors` | PASS | Successfully aggregates multiple detectors within one group. |
| `TestParseGarakEvalToOCSF_MultipleGroups` | PASS | Successfully aggregates multiple groups in one run. |
| `TestParseGarakEvalToOCSF_MissingEntryType` | PASS | Rejects lines missing `entry_type`. |
| `TestParseGarakEvalToOCSF_ProductMetadata` | PASS | Hardcodes Product as `SAE Tool (Garak Scanner)`. |
| `TestParseGarakEvalToOCSF_ActivityName` | PASS | Hardcodes ActivityName as `AI Vulnerability Scan`. |
| `TestParseGarakEvalToOCSF_BadEvalStructure` | PASS | Handles schema drifts or nested string mappings cleanly. |

## 2. Regression Baseline
Existing SAE integrations were evaluated in tandem. No side effects detected.
- **Wazuh**: PASS (4/4)
- **Zeek**: PASS (4/4)
- **Suricata**: PASS (7/7)
- **Falco**: PASS (12/12)
- **KubeArmor**: PASS (14/14)
- **Trivy**: PASS (18/18)
- **ScoutSuite**: PASS (18/18)
- **Shuffle**: PASS (18/18)
- **TheHive**: PASS (18/18)
- **Cortex**: PASS (18/18)
- **Ollama**: PASS (19/19)
- **LangGraph**: PASS (10/10)
- **Garak**: PASS (18/18)

**Total Test Count**: 178/178 passing.
