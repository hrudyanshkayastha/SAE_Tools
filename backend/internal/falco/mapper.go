package falco

import (
	"encoding/json"
	"fmt"
	"sae-core/models"
)

// MapFalcoToOCSF maps a raw Falco JSON event to the unified OCSF format.
func MapFalcoToOCSF(rawJSON []byte) (*models.OCSFFinding, error) {
	var falcoEvent FalcoEvent
	if err := json.Unmarshal(rawJSON, &falcoEvent); err != nil {
		return nil, fmt.Errorf("failed to parse Falco record: %w", err)
	}

	// Validate required fields
	if falcoEvent.Output == "" && falcoEvent.Rule == "" {
		return nil, fmt.Errorf("missing required Falco fields (rule/output)")
	}

	finding := &models.OCSFFinding{
		Time:         falcoEvent.Time,
		ActivityName: "Runtime Security Event",
	}

	finding.Metadata.Product = "SAE Tool (Falco Engine)"
	finding.Metadata.Version = "1.0"
	finding.Metadata.LogFormat = "Falco JSON"

	finding.Message = fmt.Sprintf("Falco Rule: %s | %s", falcoEvent.Rule, falcoEvent.Output)
	finding.Severity = mapPriorityToSeverity(falcoEvent.Priority)

	// Preserve Context from OutputFields
	if containerID, ok := extractString(falcoEvent.OutputFields, "container.id"); ok && containerID != "" && containerID != "host" {
		finding.Observables = append(finding.Observables, models.Observable{
			Type:  "Container",
			Value: containerID,
		})
	}
	
	if user, ok := extractString(falcoEvent.OutputFields, "user.name"); ok && user != "" {
		finding.Observables = append(finding.Observables, models.Observable{
			Type:  "User",
			Value: user,
		})
	}

	if procCmdline, ok := extractString(falcoEvent.OutputFields, "proc.cmdline"); ok && procCmdline != "" {
		finding.Observables = append(finding.Observables, models.Observable{
			Type:  "Process",
			Value: procCmdline,
		})
	}

	if fdName, ok := extractString(falcoEvent.OutputFields, "fd.name"); ok && fdName != "" {
		finding.Observables = append(finding.Observables, models.Observable{
			Type:  "File",
			Value: fdName,
		})
	}

	return finding, nil
}

func mapPriorityToSeverity(priority string) string {
	switch priority {
	case "Emergency":
		return "Fatal"
	case "Alert":
		return "Critical"
	case "Critical":
		return "Critical"
	case "Error":
		return "High"
	case "Warning":
		return "Medium"
	case "Notice":
		return "Low"
	case "Informational":
		return "Info"
	case "Debug":
		return "Info"
	default:
		return "Unknown"
	}
}

func extractString(fields map[string]interface{}, key string) (string, bool) {
	if val, exists := fields[key]; exists {
		str, ok := val.(string)
		return str, ok
	}
	return "", false
}
