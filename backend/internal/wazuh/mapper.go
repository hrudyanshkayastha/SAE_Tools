package wazuh

import (
	"encoding/json"
	"fmt"
	"sae-core/models"
)

// MapAlertToOCSF converts a native Wazuh JSON alert into the SAE OCSF format.
// This abstract mapping ensures SAE is decoupled from Wazuh's specific syntax.
func MapAlertToOCSF(rawJSON []byte) (*models.OCSFFinding, error) {
	var alert Alert
	if err := json.Unmarshal(rawJSON, &alert); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wazuh alert: %w", err)
	}

	finding := &models.OCSFFinding{
		ActivityID:   1, // 1 = New Alert
		ActivityName: "Alert",
		Time:         alert.Timestamp,
		Message:      alert.Rule.Description,
		Status:       "New",
		RawEvent:     string(rawJSON),
	}

	// Map Wazuh Level (1-15) to OCSF Severity
	// 1-3 = Info, 4-6 = Low, 7-9 = Medium, 10-12 = High, 13-14 = Critical, 15 = Fatal
	switch {
	case alert.Rule.Level <= 3:
		finding.SeverityID = 1
		finding.Severity = "Info"
	case alert.Rule.Level <= 6:
		finding.SeverityID = 2
		finding.Severity = "Low"
	case alert.Rule.Level <= 9:
		finding.SeverityID = 3
		finding.Severity = "Medium"
	case alert.Rule.Level <= 12:
		finding.SeverityID = 4
		finding.Severity = "High"
	case alert.Rule.Level <= 14:
		finding.SeverityID = 5
		finding.Severity = "Critical"
	case alert.Rule.Level >= 15:
		finding.SeverityID = 6
		finding.Severity = "Fatal"
	}

	// Set Metadata
	finding.Metadata.Product = "SAE Tool (Wazuh Engine)"
	finding.Metadata.Vendor = "SAE"
	finding.Metadata.LogFormat = "wazuh_json"

	// Extract Observables
	if alert.Data.SrcIP != "" {
		finding.Observables = append(finding.Observables, models.Observable{
			Name:  "Source IP",
			Type:  "IP",
			Value: alert.Data.SrcIP,
		})
	}
	if alert.Data.DstIP != "" {
		finding.Observables = append(finding.Observables, models.Observable{
			Name:  "Destination IP",
			Type:  "IP",
			Value: alert.Data.DstIP,
		})
	}
	if alert.Agent.IP != "" {
		finding.Observables = append(finding.Observables, models.Observable{
			Name:  "Agent IP",
			Type:  "IP",
			Value: alert.Agent.IP,
		})
	}

	return finding, nil
}
