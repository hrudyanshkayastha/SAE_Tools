package shuffle

import (
	"strings"
	"testing"
)

// SYNTHETIC FIXTURE — UNIT TEST ONLY
func TestMapResultToOCSF_Valid(t *testing.T) {
	raw := `{"execution_id":"123","status":"SUCCESS","action":{"app_name":"hello","name":"world","label":"lbl"},"result":"done"}`
	f, err := MapResultToOCSF([]byte(raw))
	if err != nil || len(f) != 1 { t.Fatalf("expected 1 finding") }
}

func TestMapResultToOCSF_StatusSuccess(t *testing.T) {
	raw := `{"execution_id":"123","status":"SUCCESS","action":{"app_name":"hello"}}`
	f, _ := MapResultToOCSF([]byte(raw))
	if f[0].Severity != "Info" { t.Errorf("expected Info") }
}

func TestMapResultToOCSF_StatusFailure(t *testing.T) {
	raw := `{"execution_id":"123","status":"FAILURE","action":{"app_name":"hello"}}`
	f, _ := MapResultToOCSF([]byte(raw))
	if f[0].Severity != "High" { t.Errorf("expected High") }
}

func TestMapResultToOCSF_StatusFailed(t *testing.T) {
	raw := `{"execution_id":"123","status":"FAILED","action":{"app_name":"hello"}}`
	f, _ := MapResultToOCSF([]byte(raw))
	if f[0].Severity != "High" { t.Errorf("expected High") }
}

func TestMapResultToOCSF_StatusError(t *testing.T) {
	raw := `{"execution_id":"123","status":"ERROR","action":{"app_name":"hello"}}`
	f, _ := MapResultToOCSF([]byte(raw))
	if f[0].Severity != "High" { t.Errorf("expected High") }
}

func TestMapResultToOCSF_StatusAborted(t *testing.T) {
	raw := `{"execution_id":"123","status":"ABORTED","action":{"app_name":"hello"}}`
	f, _ := MapResultToOCSF([]byte(raw))
	if f[0].Severity != "Medium" { t.Errorf("expected Medium") }
}

func TestMapResultToOCSF_StatusUnknown(t *testing.T) {
	raw := `{"execution_id":"123","status":"PENDING","action":{"app_name":"hello"}}`
	f, _ := MapResultToOCSF([]byte(raw))
	if f[0].Severity != "Info" { t.Errorf("expected Info") }
}

func TestMapResultToOCSF_ExecutionIDObservable(t *testing.T) {
	raw := `{"execution_id":"exec-123","action":{"app_name":"hello"}}`
	f, _ := MapResultToOCSF([]byte(raw))
	found := false
	for _, o := range f[0].Observables { if o.Type == "ExecutionID" && o.Value == "exec-123" { found = true } }
	if !found { t.Errorf("missing execution id observable") }
}

func TestMapResultToOCSF_AppNameObservable(t *testing.T) {
	raw := `{"execution_id":"123","action":{"app_name":"my-app"}}`
	f, _ := MapResultToOCSF([]byte(raw))
	found := false
	for _, o := range f[0].Observables { if o.Type == "AppName" && o.Value == "my-app" { found = true } }
	if !found { t.Errorf("missing app name observable") }
}

func TestMapResultToOCSF_ActionNameObservable(t *testing.T) {
	raw := `{"execution_id":"123","action":{"app_name":"my-app","name":"my-action"}}`
	f, _ := MapResultToOCSF([]byte(raw))
	found := false
	for _, o := range f[0].Observables { if o.Type == "ActionName" && o.Value == "my-action" { found = true } }
	if !found { t.Errorf("missing action name observable") }
}

func TestMapResultToOCSF_StatusObservable(t *testing.T) {
	raw := `{"execution_id":"123","status":"SUCCESS","action":{"app_name":"my-app"}}`
	f, _ := MapResultToOCSF([]byte(raw))
	found := false
	for _, o := range f[0].Observables { if o.Type == "Status" && o.Value == "SUCCESS" { found = true } }
	if !found { t.Errorf("missing status observable") }
}

func TestMapResultToOCSF_ResultTruncation(t *testing.T) {
	longStr := strings.Repeat("a", 300)
	raw := `{"execution_id":"123","action":{"app_name":"my-app"},"result":"` + longStr + `"}`
	f, _ := MapResultToOCSF([]byte(raw))
	if len(f[0].Message) > 300 { t.Errorf("result should be truncated") }
	if !strings.Contains(f[0].Message, "...") { t.Errorf("expected ... for truncated result") }
}

func TestMapResultToOCSF_MissingOptionalFields(t *testing.T) {
	raw := `{"execution_id":"123"}`
	f, err := MapResultToOCSF([]byte(raw))
	if err != nil || len(f) != 1 { t.Errorf("should handle missing optional items") }
}

func TestMapResultToOCSF_MissingRequiredFields(t *testing.T) {
	raw := `{"status":"SUCCESS"}`
	_, err := MapResultToOCSF([]byte(raw))
	if err == nil { t.Errorf("should fail empty report") }
}

func TestMapResultToOCSF_MalformedJSON(t *testing.T) {
	raw := `{"execution_id":`
	_, err := MapResultToOCSF([]byte(raw))
	if err == nil { t.Errorf("should fail malformed") }
}

func TestMapResultToOCSF_InvalidFieldTypes(t *testing.T) {
	raw := `{"execution_id": 123}`
	_, err := MapResultToOCSF([]byte(raw))
	if err == nil { t.Errorf("should fail invalid types") }
}

func TestMapResultToOCSF_EmptyAssessment(t *testing.T) {
	raw := `{}`
	_, err := MapResultToOCSF([]byte(raw))
	if err == nil { t.Errorf("should fail empty report") }
}

func TestMapResultToOCSF_ShortResultNoTruncation(t *testing.T) {
	raw := `{"execution_id":"123","action":{"app_name":"my-app"},"result":"short"}`
	f, _ := MapResultToOCSF([]byte(raw))
	if strings.Contains(f[0].Message, "...") { t.Errorf("should not truncate short result") }
}
