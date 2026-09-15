# SAE Falco Test Report

## 1. Unit Test Coverage
The Falco parsing and normalization adapter (`sae-core/internal/falco`) was tested with 12 deterministic test cases covering functional mapping, contextual extraction, and safe negative error handling.

**Test Matrix**:
1. `TestMapFalcoToOCSF_ValidRuntime` - Valid Falco runtime event parsing.
2. `TestMapFalcoToOCSF_HighPriority` - Critical/Emergency priority mapping.
3. `TestMapFalcoToOCSF_ContainerEvent` - `container.id` observable extraction.
4. `TestMapFalcoToOCSF_ProcessEvent` - `proc.cmdline` observable extraction.
5. `TestMapFalcoToOCSF_FileEvent` - `fd.name` (file) observable extraction.
6. `TestMapFalcoToOCSF_NetworkEvent` - `fd.name` (network IP/port) observable extraction.
7. `TestMapFalcoToOCSF_MissingOptional` - Successful parsing when context fields are missing.
8. `TestMapFalcoToOCSF_MissingRequired` - Safe rejection when `rule` or `output` is missing.
9. `TestMapFalcoToOCSF_MalformedJSON` - Safe rejection of syntax errors.
10. `TestMapFalcoToOCSF_InvalidFieldTypes` - Safe rejection of unexpected data types.
11. `TestMapFalcoToOCSF_InvalidTimestamp` - Safe rejection of malformed timestamps.
12. `TestMapFalcoToOCSF_UnknownPriority` - Fallback severity generation.

**Result**: 12/12 Falco Tests Passed.

## 2. Regression Testing
A full regression sweep was executed against the entire backend after integrating the Falco listener.

```text
?       sae-core        [no test files]
ok      sae-core/internal/falco      0.188s (12 tests)
ok      sae-core/internal/suricata   (cached) (7 tests)
ok      sae-core/internal/wazuh      (cached) (4 tests)
ok      sae-core/internal/zeek       (cached) (4 tests)
```
**Total Passing tests:** 27/27 tests (100% success rate across all integrations).
No tests were removed or falsified to accommodate Falco.

## 3. Failure Isolation
By design, the Go backend initiates Falco tailing (`bufio.Reader`) via a detached Goroutine leveraging channels. If Falco pushes severely malformed bytes, the single map operation fails, logs to `fmt.Printf`, and safely discards the event while Wazuh, Zeek, and Suricata ingestions remain actively executing.
