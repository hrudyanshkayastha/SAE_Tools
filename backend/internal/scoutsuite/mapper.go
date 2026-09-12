package scoutsuite

import (
	"encoding/json"
	"fmt"
	"sae-core/models"
	"time"
)

func mapSeverity(level string) string {
	switch level {
	case "danger":
		return "Critical"
	case "warning":
		return "High"
	case "info":
		return "Medium"
	default:
		return "Info"
	}
}

func MapReportToOCSF(rawEvent []byte) ([]models.OCSFFinding, error) {
	var report ScoutSuiteReport
	if err := json.Unmarshal(rawEvent, &report); err != nil {
		return nil, fmt.Errorf("failed to parse ScoutSuite report: %w", err)
	}

	if report.ProviderCode == "" && len(report.Services) == 0 {
		return nil, fmt.Errorf("empty or invalid ScoutSuite report")
	}

	var findings []models.OCSFFinding

	for svcName, svc := range report.Services {
		for checkId, finding := range svc.Findings {
			// Skip empty findings
			if finding.Title == "" && len(finding.Items) == 0 {
				continue
			}

			// Base finding
			ocsf := models.OCSFFinding{
				Time:         time.Now(),
				ActivityName: "Cloud Assessment",
				Severity:     mapSeverity(finding.Level),
			}
			
			ocsf.Metadata.Product = "SAE Tool (ScoutSuite Engine)"

			ocsf.Message = fmt.Sprintf("ScoutSuite Finding: [%s] %s (Service: %s)", checkId, finding.Title, svcName)
			if finding.Remediation != "" {
				ocsf.Message += fmt.Sprintf(" | Remediation: %s", finding.Remediation)
			}

			// Add observables for context
			ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "CloudProvider", Value: report.ProviderCode})
			if report.AccountID != "" {
				ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "CloudAccount", Value: report.AccountID})
			}
			ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "CloudService", Value: svcName})
			ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "CheckID", Value: checkId})

			// Add an observable for each flagged item (resource identifier)
			for _, item := range finding.Items {
				ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "CloudResource", Value: item})
			}

			// In OCSF, a single check could apply to multiple items.
			// But for simplicity in the finding log, we emit one aggregated finding per check.
			findings = append(findings, ocsf)
		}
	}

	return findings, nil
}
