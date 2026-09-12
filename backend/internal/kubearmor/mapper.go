package kubearmor

import (
	"encoding/json"
	"fmt"
	"sae-core/models"
	"time"
)

func mapSeverity(severity string) string {
	switch severity {
	case "1":
		return "Fatal"
	case "2", "3":
		return "Critical"
	case "4", "5", "6":
		return "High"
	case "7", "8":
		return "Medium"
	case "9", "10":
		return "Low"
	default:
		return "Info"
	}
}

func MapAlertToOCSF(rawEvent []byte) (*models.OCSFFinding, error) {
	var alert KubeArmorAlert
	if err := json.Unmarshal(rawEvent, &alert); err != nil {
		return nil, fmt.Errorf("failed to parse KubeArmor event: %w", err)
	}

	if alert.PolicyName == "" && alert.Operation == "" {
		return nil, fmt.Errorf("missing required KubeArmor fields")
	}

	finding := &models.OCSFFinding{
		Time:         time.Now(),
		ActivityName: "Runtime Security Event",
		Severity:     mapSeverity(alert.Severity),
		Message:      fmt.Sprintf("KubeArmor %s: %s on %s", alert.Action, alert.Operation, alert.Resource),
	}
	finding.Metadata.Product = "SAE Tool (KubeArmor Engine)"

	if alert.PolicyName != "" {
		finding.Message = fmt.Sprintf("KubeArmor Policy Violation: %s | %s", alert.PolicyName, finding.Message)
	}
	if alert.Result != "" {
		finding.Message += fmt.Sprintf(" | Result: %s", alert.Result)
	}

	if alert.ContainerID != "" {
		finding.Observables = append(finding.Observables, models.Observable{Type: "Container", Value: alert.ContainerID})
	}
	if alert.ProcessName != "" {
		finding.Observables = append(finding.Observables, models.Observable{Type: "Process", Value: alert.ProcessName})
	}
	if alert.Resource != "" {
		finding.Observables = append(finding.Observables, models.Observable{Type: "File", Value: alert.Resource})
	}
	if alert.NamespaceName != "" {
		finding.Observables = append(finding.Observables, models.Observable{Type: "Namespace", Value: alert.NamespaceName})
	}
	if alert.PodName != "" {
		finding.Observables = append(finding.Observables, models.Observable{Type: "Pod", Value: alert.PodName})
	}

	return finding, nil
}
