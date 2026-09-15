package trivy

import (
	"encoding/json"
	"fmt"
	"sae-core/models"
	"time"
)

func mapSeverity(severity string) string {
	switch severity {
	case "CRITICAL":
		return "Fatal"
	case "HIGH":
		return "Critical"
	case "MEDIUM":
		return "High"
	case "LOW":
		return "Medium"
	case "UNKNOWN":
		return "Info"
	default:
		return "Info"
	}
}

func MapReportToOCSF(rawEvent []byte) ([]models.OCSFFinding, error) {
	var report TrivyReport
	if err := json.Unmarshal(rawEvent, &report); err != nil {
		return nil, fmt.Errorf("failed to parse Trivy report: %w", err)
	}

	var findings []models.OCSFFinding

	for _, result := range report.Results {
		for _, vuln := range result.Vulnerabilities {
			if vuln.VulnerabilityID == "" {
				continue
			}

			finding := models.OCSFFinding{
				Time:         time.Now(),
				ActivityName: "Vulnerability Discovery",
				Severity:     mapSeverity(vuln.Severity),
				Message:      fmt.Sprintf("Trivy Vulnerability: %s (%s) in %s %s", vuln.VulnerabilityID, vuln.Severity, vuln.PkgName, vuln.InstalledVersion),
			}
			finding.Metadata.Product = "SAE Tool (Trivy Engine)"

			if vuln.FixedVersion != "" {
				finding.Message += fmt.Sprintf(" | Fixed in: %s", vuln.FixedVersion)
			}
			if vuln.Title != "" {
				finding.Message += fmt.Sprintf(" | %s", vuln.Title)
			}

			finding.Observables = append(finding.Observables, models.Observable{Type: "Vulnerability", Value: vuln.VulnerabilityID})
			finding.Observables = append(finding.Observables, models.Observable{Type: "Package", Value: vuln.PkgName})
			
			if report.ArtifactName != "" {
				finding.Observables = append(finding.Observables, models.Observable{Type: "Artifact", Value: report.ArtifactName})
			}
			if result.Target != "" {
				finding.Observables = append(finding.Observables, models.Observable{Type: "Target", Value: result.Target})
			}

			findings = append(findings, finding)
		}
	}

	if len(findings) == 0 && len(report.Results) == 0 && report.ArtifactName == "" {
		return nil, fmt.Errorf("empty or invalid trivy report")
	}

	return findings, nil
}
