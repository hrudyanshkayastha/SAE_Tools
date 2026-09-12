package scoutsuite

import (
	"strings"
	"testing"
)

// SYNTHETIC FIXTURE — UNIT TEST ONLY
func TestMapReportToOCSF_ValidFinding(t *testing.T) {
	raw := `{"provider_code":"aws","services":{"ec2":{"findings":{"check-1":{"title":"Test","level":"danger","items":["res1"]}}}}}`
	f, err := MapReportToOCSF([]byte(raw))
	if err != nil || len(f) != 1 { t.Fatalf("expected 1 finding") }
}

func TestMapReportToOCSF_SeverityHigh(t *testing.T) {
	raw := `{"provider_code":"aws","services":{"s3":{"findings":{"check-1":{"title":"Test","level":"danger"}}}}}`
	f, _ := MapReportToOCSF([]byte(raw))
	if f[0].Severity != "Critical" { t.Errorf("expected Critical") }
}

func TestMapReportToOCSF_SeverityMedium(t *testing.T) {
	raw := `{"provider_code":"aws","services":{"s3":{"findings":{"check-1":{"title":"Test","level":"warning"}}}}}`
	f, _ := MapReportToOCSF([]byte(raw))
	if f[0].Severity != "High" { t.Errorf("expected High") }
}

func TestMapReportToOCSF_SeverityLow(t *testing.T) {
	raw := `{"provider_code":"aws","services":{"s3":{"findings":{"check-1":{"title":"Test","level":"info"}}}}}`
	f, _ := MapReportToOCSF([]byte(raw))
	if f[0].Severity != "Medium" { t.Errorf("expected Medium") }
}

func TestMapReportToOCSF_ProviderMetadata(t *testing.T) {
	raw := `{"provider_code":"gcp","services":{"compute":{"findings":{"check-1":{"title":"Test"}}}}}`
	f, _ := MapReportToOCSF([]byte(raw))
	found := false
	for _, o := range f[0].Observables { if o.Type == "CloudProvider" && o.Value == "gcp" { found = true } }
	if !found { t.Errorf("missing provider") }
}

func TestMapReportToOCSF_AccountMetadata(t *testing.T) {
	raw := `{"provider_code":"aws","account_id":"123","services":{"s3":{"findings":{"check-1":{"title":"Test"}}}}}`
	f, _ := MapReportToOCSF([]byte(raw))
	found := false
	for _, o := range f[0].Observables { if o.Type == "CloudAccount" && o.Value == "123" { found = true } }
	if !found { t.Errorf("missing account id") }
}

func TestMapReportToOCSF_ServiceMetadata(t *testing.T) {
	raw := `{"provider_code":"aws","services":{"iam":{"findings":{"check-1":{"title":"Test"}}}}}`
	f, _ := MapReportToOCSF([]byte(raw))
	found := false
	for _, o := range f[0].Observables { if o.Type == "CloudService" && o.Value == "iam" { found = true } }
	if !found { t.Errorf("missing service") }
}

func TestMapReportToOCSF_ResourceMetadata(t *testing.T) {
	raw := `{"provider_code":"aws","services":{"iam":{"findings":{"check-1":{"title":"Test","items":["user1"]}}}}}`
	f, _ := MapReportToOCSF([]byte(raw))
	found := false
	for _, o := range f[0].Observables { if o.Type == "CloudResource" && o.Value == "user1" { found = true } }
	if !found { t.Errorf("missing resource") }
}

func TestMapReportToOCSF_CheckIdentifier(t *testing.T) {
	raw := `{"provider_code":"aws","services":{"iam":{"findings":{"my-check":{"title":"Test"}}}}}`
	f, _ := MapReportToOCSF([]byte(raw))
	if !strings.Contains(f[0].Message, "[my-check]") { t.Errorf("missing check id") }
}

func TestMapReportToOCSF_Remediation(t *testing.T) {
	raw := `{"provider_code":"aws","services":{"iam":{"findings":{"check-1":{"title":"Test","remediation":"Fix it"}}}}}`
	f, _ := MapReportToOCSF([]byte(raw))
	if !strings.Contains(f[0].Message, "Fix it") { t.Errorf("missing remediation") }
}

func TestMapReportToOCSF_MultipleFindings(t *testing.T) {
	raw := `{"provider_code":"aws","services":{"iam":{"findings":{"check-1":{"title":"A"},"check-2":{"title":"B"}}}}}`
	f, _ := MapReportToOCSF([]byte(raw))
	if len(f) != 2 { t.Errorf("expected 2 findings") }
}

func TestMapReportToOCSF_MissingOptional(t *testing.T) {
	raw := `{"provider_code":"aws","services":{"iam":{"findings":{"check-1":{"title":"A"}}}}}`
	f, err := MapReportToOCSF([]byte(raw))
	if err != nil || len(f) != 1 { t.Errorf("should handle missing optional items array") }
}

func TestMapReportToOCSF_MissingRequired(t *testing.T) {
	raw := `{"account_id":"123"}`
	_, err := MapReportToOCSF([]byte(raw))
	if err == nil { t.Errorf("should fail empty report") }
}

func TestMapReportToOCSF_MalformedJSON(t *testing.T) {
	raw := `{"provider_code":`
	_, err := MapReportToOCSF([]byte(raw))
	if err == nil { t.Errorf("should fail malformed") }
}

func TestMapReportToOCSF_InvalidFieldTypes(t *testing.T) {
	raw := `{"provider_code": 123}`
	_, err := MapReportToOCSF([]byte(raw))
	if err == nil { t.Errorf("should fail invalid types") }
}

func TestMapReportToOCSF_UnknownCheckType(t *testing.T) {
	raw := `{"provider_code":"aws","services":{"iam":{"findings":{"check-1":{"title":"Test","level":"super-critical"}}}}}`
	f, _ := MapReportToOCSF([]byte(raw))
	if f[0].Severity != "Info" { t.Errorf("should default to Info") }
}

func TestMapReportToOCSF_EmptyAssessment(t *testing.T) {
	raw := `{"provider_code":"aws","services":{}}`
	f, err := MapReportToOCSF([]byte(raw))
	if err != nil || len(f) != 0 { t.Errorf("should parse empty without error") }
}

func TestMapReportToOCSF_DuplicateItems(t *testing.T) {
	raw := `{"provider_code":"aws","services":{"iam":{"findings":{"check-1":{"title":"Test","items":["res1", "res2"]}}}}}`
	f, _ := MapReportToOCSF([]byte(raw))
	count := 0
	for _, o := range f[0].Observables { if o.Type == "CloudResource" { count++ } }
	if count != 2 { t.Errorf("should have 2 resources") }
}
