package suricata

import (
	"encoding/json"
	"fmt"
	"sae-core/models"
)

// MapEVEToOCSF maps a raw Suricata EVE JSON event to the unified OCSF format.
func MapEVEToOCSF(rawJSON []byte) (*models.OCSFFinding, error) {
	var eve SuricataEVE
	if err := json.Unmarshal(rawJSON, &eve); err != nil {
		return nil, fmt.Errorf("failed to parse Suricata EVE record: %w", err)
	}

	finding := &models.OCSFFinding{
		Time: eve.Timestamp,
	}
	finding.Metadata.Product = "SAE Tool (Suricata Engine)"
	finding.Metadata.Version = "1.0"

	finding.Observables = append(finding.Observables, models.Observable{
		Type:  "IP",
		Value: eve.SrcIP,
	})
	finding.Observables = append(finding.Observables, models.Observable{
		Type:  "IP",
		Value: eve.DestIP,
	})

	switch eve.EventType {
	case "alert":
		finding.ActivityName = "Intrusion Detection"
		if eve.Alert != nil {
			finding.Message = fmt.Sprintf("Suricata Alert: %s (Category: %s)", eve.Alert.Signature, eve.Alert.Category)
			finding.Severity = mapSeverity(eve.Alert.Severity)
		} else {
			finding.Message = "Suricata Alert (Missing Alert Details)"
			finding.Severity = "High"
		}
	case "flow":
		finding.ActivityName = "Network Activity"
		finding.Message = fmt.Sprintf("Network flow from %s:%d to %s:%d over %s", eve.SrcIP, eve.SrcPort, eve.DestIP, eve.DestPort, eve.Proto)
		finding.Severity = "Info"
	case "dns":
		finding.ActivityName = "DNS Activity"
		if eve.DNS != nil {
			finding.Message = fmt.Sprintf("DNS %s for %s (%s)", eve.DNS.Type, eve.DNS.RRName, eve.DNS.RRType)
		} else {
			finding.Message = "DNS Activity"
		}
		finding.Severity = "Info"
	case "http":
		finding.ActivityName = "HTTP Activity"
		if eve.HTTP != nil {
			finding.Message = fmt.Sprintf("HTTP %s %s%s", eve.HTTP.HTTPMethod, eve.HTTP.Hostname, eve.HTTP.URL)
		} else {
			finding.Message = "HTTP Activity"
		}
		finding.Severity = "Info"
	case "tls":
		finding.ActivityName = "TLS Activity"
		if eve.TLS != nil {
			finding.Message = fmt.Sprintf("TLS Connection SNI: %s, Version: %s", eve.TLS.Sni, eve.TLS.Version)
		} else {
			finding.Message = "TLS Activity"
		}
		finding.Severity = "Info"
	default:
		finding.ActivityName = "Network Telemetry"
		finding.Message = fmt.Sprintf("Suricata Event [%s] from %s to %s", eve.EventType, eve.SrcIP, eve.DestIP)
		finding.Severity = "Info"
	}

	return finding, nil
}

func mapSeverity(suricataSeverity int) string {
	// Suricata: 1 is high/alert, 2 is medium, 3 is low, 4 is info
	switch suricataSeverity {
	case 1:
		return "Critical"
	case 2:
		return "High"
	case 3:
		return "Medium"
	case 4:
		return "Low"
	default:
		return "Unknown"
	}
}
