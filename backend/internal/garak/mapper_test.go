package garak

import (
	"strings"
	"testing"
)

func TestParseGarakEvalToOCSF_Empty(t *testing.T) {
	_, err := ParseGarakEvalToOCSF([]byte(""))
	if err == nil { t.Errorf("expected error on empty") }
}

func TestParseGarakEvalToOCSF_InvalidJSON(t *testing.T) {
	_, err := ParseGarakEvalToOCSF([]byte("{invalid"))
	if err == nil { t.Errorf("expected error on invalid json") }
}

func TestParseGarakEvalToOCSF_NotEval(t *testing.T) {
	res, err := ParseGarakEvalToOCSF([]byte(`{"entry_type": "start_run setup"}`))
	if err != nil { t.Fatalf("should not error") }
	if len(res) != 0 { t.Errorf("should ignore non-eval lines") }
}

func TestParseGarakEvalToOCSF_ValidEval(t *testing.T) {
	jsonl := `{"entry_type":"eval","eval":{"lmrc":{"_summary":{},"probe1":{"_summary":{}},"det1":{"detector_name":"det1","detector_descr":"test","passed":2,"total_evaluated":2,"absolute_defcon":5}}}}`
	f, err := ParseGarakEvalToOCSF([]byte(jsonl))
	if err != nil { t.Fatalf("failed: %v", err) }
	if len(f) != 1 { t.Fatalf("expected 1 finding") }
	if f[0].Severity != "Info" { t.Errorf("expected Info since passed == total") }
}

func TestParseGarakEvalToOCSF_SeverityCritical(t *testing.T) {
	jsonl := `{"entry_type":"eval","eval":{"lmrc":{"det1":{"detector_name":"det1","passed":0,"total_evaluated":2,"absolute_defcon":2}}}}`
	f, _ := ParseGarakEvalToOCSF([]byte(jsonl))
	if f[0].Severity != "Critical" { t.Errorf("expected Critical") }
}

func TestParseGarakEvalToOCSF_SeverityHigh(t *testing.T) {
	jsonl := `{"entry_type":"eval","eval":{"lmrc":{"det1":{"detector_name":"det1","passed":1,"total_evaluated":2,"absolute_defcon":3}}}}`
	f, _ := ParseGarakEvalToOCSF([]byte(jsonl))
	if f[0].Severity != "High" { t.Errorf("expected High") }
}

func TestParseGarakEvalToOCSF_SeverityMedium(t *testing.T) {
	jsonl := `{"entry_type":"eval","eval":{"lmrc":{"det1":{"detector_name":"det1","passed":1,"total_evaluated":2,"absolute_defcon":4}}}}`
	f, _ := ParseGarakEvalToOCSF([]byte(jsonl))
	if f[0].Severity != "Medium" { t.Errorf("expected Medium") }
}

func TestParseGarakEvalToOCSF_SeverityLow(t *testing.T) {
	jsonl := `{"entry_type":"eval","eval":{"lmrc":{"det1":{"detector_name":"det1","passed":1,"total_evaluated":2,"absolute_defcon":5}}}}`
	f, _ := ParseGarakEvalToOCSF([]byte(jsonl))
	if f[0].Severity != "Low" { t.Errorf("expected Low") }
}

func TestParseGarakEvalToOCSF_Observables(t *testing.T) {
	jsonl := `{"entry_type":"eval","eval":{"group1":{"det1":{"detector_name":"det1","passed":1,"total_evaluated":2}}}}`
	f, _ := ParseGarakEvalToOCSF([]byte(jsonl))
	foundDet := false
	foundGrp := false
	for _, o := range f[0].Observables {
		if o.Type == "Detector" && o.Value == "det1" { foundDet = true }
		if o.Type == "Group" && o.Value == "group1" { foundGrp = true }
	}
	if !foundDet || !foundGrp { t.Errorf("missing observables") }
}

func TestParseGarakEvalToOCSF_IgnoresSummary(t *testing.T) {
	jsonl := `{"entry_type":"eval","eval":{"group1":{"_summary":{"detector_name":"fake"}}}}`
	f, _ := ParseGarakEvalToOCSF([]byte(jsonl))
	if len(f) != 0 { t.Errorf("should ignore _summary") }
}

func TestParseGarakEvalToOCSF_IgnoresProbes(t *testing.T) {
	jsonl := `{"entry_type":"eval","eval":{"group1":{"probe1":{"_summary":{}}}}}`
	f, _ := ParseGarakEvalToOCSF([]byte(jsonl))
	if len(f) != 0 { t.Errorf("should ignore probes without detector_name") }
}

func TestParseGarakEvalToOCSF_MessageFormat(t *testing.T) {
	jsonl := `{"entry_type":"eval","eval":{"group1":{"det1":{"detector_name":"det1","detector_descr":"A test","passed":1,"total_evaluated":5}}}}`
	f, _ := ParseGarakEvalToOCSF([]byte(jsonl))
	if !strings.Contains(f[0].Message, "Garak Detector [det1]: A test (Passed 1/5)") {
		t.Errorf("incorrect message format: %s", f[0].Message)
	}
}

func TestParseGarakEvalToOCSF_MultipleDetectors(t *testing.T) {
	jsonl := `{"entry_type":"eval","eval":{"group1":{"det1":{"detector_name":"det1"},"det2":{"detector_name":"det2"}}}}`
	f, _ := ParseGarakEvalToOCSF([]byte(jsonl))
	if len(f) != 2 { t.Errorf("expected 2 findings") }
}

func TestParseGarakEvalToOCSF_MultipleGroups(t *testing.T) {
	jsonl := `{"entry_type":"eval","eval":{"group1":{"det1":{"detector_name":"det1"}},"group2":{"det2":{"detector_name":"det2"}}}}`
	f, _ := ParseGarakEvalToOCSF([]byte(jsonl))
	if len(f) != 2 { t.Errorf("expected 2 findings across groups") }
}

func TestParseGarakEvalToOCSF_MissingEntryType(t *testing.T) {
	jsonl := `{"eval":{}}`
	f, _ := ParseGarakEvalToOCSF([]byte(jsonl))
	if len(f) != 0 { t.Errorf("should ignore") }
}

func TestParseGarakEvalToOCSF_ProductMetadata(t *testing.T) {
	jsonl := `{"entry_type":"eval","eval":{"group1":{"det1":{"detector_name":"det1"}}}}`
	f, _ := ParseGarakEvalToOCSF([]byte(jsonl))
	if f[0].Metadata.Product != "SAE Tool (Garak Scanner)" { t.Errorf("bad product") }
}

func TestParseGarakEvalToOCSF_ActivityName(t *testing.T) {
	jsonl := `{"entry_type":"eval","eval":{"group1":{"det1":{"detector_name":"det1"}}}}`
	f, _ := ParseGarakEvalToOCSF([]byte(jsonl))
	if f[0].ActivityName != "AI Vulnerability Scan" { t.Errorf("bad activity name") }
}

func TestParseGarakEvalToOCSF_BadEvalStructure(t *testing.T) {
	jsonl := `{"entry_type":"eval","eval":"not_a_map"}`
	_, err := ParseGarakEvalToOCSF([]byte(jsonl))
	if err == nil { t.Errorf("expected error mapping eval") }
}
