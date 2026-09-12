package trivy

import (
	"strings"
	"testing"
)

// SYNTHETIC TRIVY FIXTURE — UNIT TEST ONLY
func TestMapReportToOCSF_ValidResult(t *testing.T) {
	raw := `{"SchemaVersion":2,"ArtifactName":"alpine","Results":[{"Target":"alpine","Vulnerabilities":[{"VulnerabilityID":"CVE-2019-14697","PkgName":"musl","InstalledVersion":"1.1.22","FixedVersion":"1.1.23","Severity":"CRITICAL","Title":"musl libc issue"}]}]}`
	findings, err := MapReportToOCSF([]byte(raw))
	if err != nil || len(findings) != 1 {
		t.Fatalf("expected 1 finding, got err: %v", err)
	}
}

func TestMapReportToOCSF_SeverityCritical(t *testing.T) {
	raw := `{"SchemaVersion":2,"Results":[{"Vulnerabilities":[{"VulnerabilityID":"CVE-1","Severity":"CRITICAL"}]}]}`
	f, _ := MapReportToOCSF([]byte(raw))
	if f[0].Severity != "Fatal" { t.Errorf("expected Fatal") }
}

func TestMapReportToOCSF_SeverityHigh(t *testing.T) {
	raw := `{"SchemaVersion":2,"Results":[{"Vulnerabilities":[{"VulnerabilityID":"CVE-1","Severity":"HIGH"}]}]}`
	f, _ := MapReportToOCSF([]byte(raw))
	if f[0].Severity != "Critical" { t.Errorf("expected Critical") }
}

func TestMapReportToOCSF_SeverityMedium(t *testing.T) {
	raw := `{"SchemaVersion":2,"Results":[{"Vulnerabilities":[{"VulnerabilityID":"CVE-1","Severity":"MEDIUM"}]}]}`
	f, _ := MapReportToOCSF([]byte(raw))
	if f[0].Severity != "High" { t.Errorf("expected High") }
}

func TestMapReportToOCSF_SeverityLow(t *testing.T) {
	raw := `{"SchemaVersion":2,"Results":[{"Vulnerabilities":[{"VulnerabilityID":"CVE-1","Severity":"LOW"}]}]}`
	f, _ := MapReportToOCSF([]byte(raw))
	if f[0].Severity != "Medium" { t.Errorf("expected Medium") }
}

func TestMapReportToOCSF_SeverityUnknown(t *testing.T) {
	raw := `{"SchemaVersion":2,"Results":[{"Vulnerabilities":[{"VulnerabilityID":"CVE-1","Severity":"UNKNOWN"}]}]}`
	f, _ := MapReportToOCSF([]byte(raw))
	if f[0].Severity != "Info" { t.Errorf("expected Info") }
}

func TestMapReportToOCSF_PackageMetadata(t *testing.T) {
	raw := `{"SchemaVersion":2,"Results":[{"Vulnerabilities":[{"VulnerabilityID":"CVE-1","PkgName":"pkgA"}]}]}`
	f, _ := MapReportToOCSF([]byte(raw))
	if !strings.Contains(f[0].Message, "pkgA") { t.Errorf("missing package metadata") }
}

func TestMapReportToOCSF_VulnerabilityID(t *testing.T) {
	raw := `{"SchemaVersion":2,"Results":[{"Vulnerabilities":[{"VulnerabilityID":"CVE-2023-1234"}]}]}`
	f, _ := MapReportToOCSF([]byte(raw))
	if !strings.Contains(f[0].Message, "CVE-2023-1234") { t.Errorf("missing vulnerability ID") }
}

func TestMapReportToOCSF_InstalledVersion(t *testing.T) {
	raw := `{"SchemaVersion":2,"Results":[{"Vulnerabilities":[{"VulnerabilityID":"CVE-1","InstalledVersion":"1.2.3"}]}]}`
	f, _ := MapReportToOCSF([]byte(raw))
	if !strings.Contains(f[0].Message, "1.2.3") { t.Errorf("missing installed version") }
}

func TestMapReportToOCSF_FixedVersion(t *testing.T) {
	raw := `{"SchemaVersion":2,"Results":[{"Vulnerabilities":[{"VulnerabilityID":"CVE-1","FixedVersion":"1.2.4"}]}]}`
	f, _ := MapReportToOCSF([]byte(raw))
	if !strings.Contains(f[0].Message, "1.2.4") { t.Errorf("missing fixed version") }
}

func TestMapReportToOCSF_ArtifactMetadata(t *testing.T) {
	raw := `{"SchemaVersion":2,"ArtifactName":"my-artifact","Results":[{"Vulnerabilities":[{"VulnerabilityID":"CVE-1"}]}]}`
	f, _ := MapReportToOCSF([]byte(raw))
	found := false
	for _, o := range f[0].Observables { if o.Type == "Artifact" { found = true } }
	if !found { t.Errorf("missing artifact metadata") }
}

func TestMapReportToOCSF_MissingOptional(t *testing.T) {
	raw := `{"SchemaVersion":2,"Results":[{"Vulnerabilities":[{"VulnerabilityID":"CVE-X"}]}]}`
	findings, err := MapReportToOCSF([]byte(raw))
	if err != nil || len(findings) != 1 { t.Fatalf("expected 1 finding") }
}

func TestMapReportToOCSF_MissingRequired(t *testing.T) {
	raw := `{"SchemaVersion":2}`
	_, err := MapReportToOCSF([]byte(raw))
	if err == nil { t.Fatalf("expected error") }
}

func TestMapReportToOCSF_MalformedJSON(t *testing.T) {
	raw := `{"SchemaVersion":2, "Results": [`
	_, err := MapReportToOCSF([]byte(raw))
	if err == nil { t.Fatalf("expected error") }
}

func TestMapReportToOCSF_InvalidFieldTypes(t *testing.T) {
	raw := `{"SchemaVersion":"two"}`
	_, err := MapReportToOCSF([]byte(raw))
	if err == nil { t.Fatalf("expected error") }
}

func TestMapReportToOCSF_UnknownVulnerabilityType(t *testing.T) {
	raw := `{"SchemaVersion":2,"Results":[{"Vulnerabilities":[{"VulnerabilityID":"NEW-TYPE","Severity":"WEIRD"}]}]}`
	f, _ := MapReportToOCSF([]byte(raw))
	if f[0].Severity != "Info" { t.Errorf("expected fallback") }
}

func TestMapReportToOCSF_EmptyResults(t *testing.T) {
	raw := `{"SchemaVersion":2,"Results":[]}`
	_, err := MapReportToOCSF([]byte(raw))
	if err == nil { t.Fatalf("expected error") }
}

func TestMapReportToOCSF_MultipleVulnerabilities(t *testing.T) {
	raw := `{"SchemaVersion":2,"Results":[{"Vulnerabilities":[{"VulnerabilityID":"CVE-1"},{"VulnerabilityID":"CVE-2"}]}]}`
	findings, _ := MapReportToOCSF([]byte(raw))
	if len(findings) != 2 { t.Fatalf("expected 2 findings") }
}
