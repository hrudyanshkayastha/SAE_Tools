package thehive

import (
	"strings"
	"testing"
)

// SYNTHETIC FIXTURE — UNIT TEST ONLY
func TestMapCaseToOCSF_ValidFlatCase(t *testing.T) {
	raw := `{"id":"c1","number":1,"title":"title","description":"desc","severity":2,"status":"Open","assignee":"admin","tags":["t1","t2"]}`
	f, err := MapCaseToOCSF([]byte(raw))
	if err != nil || len(f) != 1 { t.Fatalf("expected 1 finding") }
}

func TestMapCaseToOCSF_ValidWebhookEvent(t *testing.T) {
	raw := `{"objectType":"case","operation":"creation","object":{"id":"c1","number":1,"title":"title"}}`
	f, err := MapCaseToOCSF([]byte(raw))
	if err != nil || len(f) != 1 { t.Fatalf("expected 1 finding") }
}

func TestMapCaseToOCSF_SeverityLow(t *testing.T) {
	raw := `{"id":"c1","title":"t","severity":1}`
	f, _ := MapCaseToOCSF([]byte(raw))
	if f[0].Severity != "Low" { t.Errorf("expected Low") }
}

func TestMapCaseToOCSF_SeverityMedium(t *testing.T) {
	raw := `{"id":"c1","title":"t","severity":2}`
	f, _ := MapCaseToOCSF([]byte(raw))
	if f[0].Severity != "Medium" { t.Errorf("expected Medium") }
}

func TestMapCaseToOCSF_SeverityHigh(t *testing.T) {
	raw := `{"id":"c1","title":"t","severity":3}`
	f, _ := MapCaseToOCSF([]byte(raw))
	if f[0].Severity != "High" { t.Errorf("expected High") }
}

func TestMapCaseToOCSF_SeverityCritical(t *testing.T) {
	raw := `{"id":"c1","title":"t","severity":4}`
	f, _ := MapCaseToOCSF([]byte(raw))
	if f[0].Severity != "Critical" { t.Errorf("expected Critical") }
}

func TestMapCaseToOCSF_SeverityUnknown(t *testing.T) {
	raw := `{"id":"c1","title":"t","severity":99}`
	f, _ := MapCaseToOCSF([]byte(raw))
	if f[0].Severity != "Info" { t.Errorf("expected Info") }
}

func TestMapCaseToOCSF_CaseIDObservable(t *testing.T) {
	raw := `{"id":"c123","title":"t"}`
	f, _ := MapCaseToOCSF([]byte(raw))
	found := false
	for _, o := range f[0].Observables { if o.Type == "CaseID" && o.Value == "c123" { found = true } }
	if !found { t.Errorf("missing case ID observable") }
}

func TestMapCaseToOCSF_StatusObservable(t *testing.T) {
	raw := `{"id":"c1","title":"t","status":"Resolved"}`
	f, _ := MapCaseToOCSF([]byte(raw))
	found := false
	for _, o := range f[0].Observables { if o.Type == "Status" && o.Value == "Resolved" { found = true } }
	if !found { t.Errorf("missing status observable") }
}

func TestMapCaseToOCSF_AssigneeObservable(t *testing.T) {
	raw := `{"id":"c1","title":"t","assignee":"john"}`
	f, _ := MapCaseToOCSF([]byte(raw))
	found := false
	for _, o := range f[0].Observables { if o.Type == "Assignee" && o.Value == "john" { found = true } }
	if !found { t.Errorf("missing assignee observable") }
}

func TestMapCaseToOCSF_TagsObservable(t *testing.T) {
	raw := `{"id":"c1","title":"t","tags":["malware","phishing"]}`
	f, _ := MapCaseToOCSF([]byte(raw))
	found := false
	for _, o := range f[0].Observables { if o.Type == "Tags" && o.Value == "malware,phishing" { found = true } }
	if !found { t.Errorf("missing tags observable") }
}

func TestMapCaseToOCSF_DescriptionTruncation(t *testing.T) {
	longStr := strings.Repeat("a", 200)
	raw := `{"id":"c1","title":"t","description":"` + longStr + `"}`
	f, _ := MapCaseToOCSF([]byte(raw))
	if len(f[0].Message) > 150 { t.Errorf("description should be truncated") }
	if !strings.Contains(f[0].Message, "...") { t.Errorf("expected ... for truncated desc") }
}

func TestMapCaseToOCSF_MissingOptionalFields(t *testing.T) {
	raw := `{"id":"c1","title":"t"}`
	f, err := MapCaseToOCSF([]byte(raw))
	if err != nil || len(f) != 1 { t.Errorf("should handle missing optional items") }
}

func TestMapCaseToOCSF_MissingRequiredFields(t *testing.T) {
	raw := `{"severity":2}`
	_, err := MapCaseToOCSF([]byte(raw))
	if err == nil { t.Errorf("should fail empty report") }
}

func TestMapCaseToOCSF_MalformedJSON(t *testing.T) {
	raw := `{"id":`
	_, err := MapCaseToOCSF([]byte(raw))
	if err == nil { t.Errorf("should fail malformed") }
}

func TestMapCaseToOCSF_InvalidFieldTypes(t *testing.T) {
	raw := `{"id": 123}`
	_, err := MapCaseToOCSF([]byte(raw))
	if err == nil { t.Errorf("should fail invalid types") }
}

func TestMapCaseToOCSF_EmptyAssessment(t *testing.T) {
	raw := `{}`
	_, err := MapCaseToOCSF([]byte(raw))
	if err == nil { t.Errorf("should fail empty report") }
}

func TestMapCaseToOCSF_ShortDescNoTruncation(t *testing.T) {
	raw := `{"id":"c1","title":"t","description":"short"}`
	f, _ := MapCaseToOCSF([]byte(raw))
	if strings.Contains(f[0].Message, "...") { t.Errorf("should not truncate short desc") }
}
