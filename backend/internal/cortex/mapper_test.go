package cortex

import (
	"strings"
	"testing"
)

func TestMapJobToOCSF_ValidJob(t *testing.T) {
	raw := `{"id":"j1","status":"Success","analyzerId":"vt","analyzerName":"VirusTotal"}`
	f, err := MapJobToOCSF([]byte(raw))
	if err != nil || len(f) != 1 { t.Fatalf("expected 1 finding") }
}

func TestMapJobToOCSF_SeveritySuccess(t *testing.T) {
	raw := `{"id":"j1","status":"Success"}`
	f, _ := MapJobToOCSF([]byte(raw))
	if f[0].Severity != "Info" { t.Errorf("expected Info") }
}

func TestMapJobToOCSF_SeverityFailure(t *testing.T) {
	raw := `{"id":"j1","status":"Failure"}`
	f, _ := MapJobToOCSF([]byte(raw))
	if f[0].Severity != "High" { t.Errorf("expected High") }
}

func TestMapJobToOCSF_MissingName(t *testing.T) {
	raw := `{"id":"j1","status":"Success","analyzerId":"vt"}`
	f, _ := MapJobToOCSF([]byte(raw))
	if !strings.Contains(f[0].Message, "vt") { t.Errorf("should fallback to ID") }
}

func TestMapJobToOCSF_Empty(t *testing.T) {
	raw := `{}`
	_, err := MapJobToOCSF([]byte(raw))
	if err == nil { t.Errorf("should fail empty object") }
}

func TestMapJobToOCSF_Malformed(t *testing.T) {
	raw := `{"id":`
	_, err := MapJobToOCSF([]byte(raw))
	if err == nil { t.Errorf("should fail malformed") }
}

func TestMapJobToOCSF_HasReport(t *testing.T) {
	raw := `{"id":"j1","report":{"summary":{"malicious":true}}}`
	f, _ := MapJobToOCSF([]byte(raw))
	found := false
	for _, o := range f[0].Observables { if o.Type == "HasReport" && o.Value == "true" { found = true } }
	if !found { t.Errorf("missing HasReport") }
}

func TestMapJobToOCSF_JobIDObservable(t *testing.T) {
	raw := `{"id":"123","analyzerId":"a1"}`
	f, _ := MapJobToOCSF([]byte(raw))
	found := false
	for _, o := range f[0].Observables { if o.Type == "JobID" && o.Value == "123" { found = true } }
	if !found { t.Errorf("missing JobID") }
}

func TestMapJobToOCSF_AnalyzerIDObservable(t *testing.T) {
	raw := `{"id":"123","analyzerId":"vt"}`
	f, _ := MapJobToOCSF([]byte(raw))
	found := false
	for _, o := range f[0].Observables { if o.Type == "AnalyzerID" && o.Value == "vt" { found = true } }
	if !found { t.Errorf("missing AnalyzerID") }
}

func TestMapJobToOCSF_InvalidDate(t *testing.T) {
	raw := `{"id":"123","date":"string"}`
	_, err := MapJobToOCSF([]byte(raw))
	if err == nil { t.Errorf("should fail type mismatch") }
}

func TestMapJobToOCSF_ValidDate(t *testing.T) {
	raw := `{"id":"123","date":1531667370000}`
	f, _ := MapJobToOCSF([]byte(raw))
	if f[0].Time.Year() < 2018 { t.Errorf("date mapping failed") }
}

func TestMapJobToOCSF_NoReport(t *testing.T) {
	raw := `{"id":"j1"}`
	f, _ := MapJobToOCSF([]byte(raw))
	for _, o := range f[0].Observables { if o.Type == "HasReport" { t.Errorf("should not have HasReport") } }
}

func TestMapJobToOCSF_EmptyReport(t *testing.T) {
	raw := `{"id":"j1","report":{}}`
	f, _ := MapJobToOCSF([]byte(raw))
	for _, o := range f[0].Observables { if o.Type == "HasReport" { t.Errorf("should not have HasReport if empty") } }
}

func TestMapJobToOCSF_StatusUnknown(t *testing.T) {
	raw := `{"id":"j1","status":"InProgress"}`
	f, _ := MapJobToOCSF([]byte(raw))
	if f[0].Severity != "Info" { t.Errorf("expected Info for InProgress") }
}

func TestMapJobToOCSF_ProductMetadata(t *testing.T) {
	raw := `{"id":"j1"}`
	f, _ := MapJobToOCSF([]byte(raw))
	if f[0].Metadata.Product != "SAE Tool (Cortex Engine)" { t.Errorf("bad product metadata") }
}

func TestMapJobToOCSF_ActivityName(t *testing.T) {
	raw := `{"id":"j1"}`
	f, _ := MapJobToOCSF([]byte(raw))
	if f[0].ActivityName != "Threat Analysis Result" { t.Errorf("bad activity name") }
}

func TestMapJobToOCSF_MissingIDFallback(t *testing.T) {
	raw := `{"analyzerId":"vt"}`
	f, _ := MapJobToOCSF([]byte(raw))
	if len(f) != 1 { t.Errorf("should process if analyzerId exists but no id") }
}

func TestMapJobToOCSF_CompletePayload(t *testing.T) {
	raw := `{"id":"j1","status":"Success","date":12345,"analyzerId":"vt","analyzerName":"VirusTotal","report":{"res":1}}`
	f, _ := MapJobToOCSF([]byte(raw))
	if len(f[0].Observables) != 3 { t.Errorf("expected 3 observables") }
}
