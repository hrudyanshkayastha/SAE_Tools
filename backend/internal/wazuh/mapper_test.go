package wazuh

import (
	"testing"
	"time"
)

func TestMapAlertToOCSF(t *testing.T) {
	rawAlert := []byte(`{
		"timestamp": "2023-10-01T12:00:00Z",
		"rule": {
			"level": 10,
			"description": "SSH Brute Force attempt",
			"id": "5712",
			"groups": ["syslog", "sshd"]
		},
		"agent": {
			"id": "001",
			"name": "web-server",
			"ip": "10.0.0.5"
		},
		"data": {
			"srcip": "192.168.1.100",
			"action": "failed"
		}
	}`)

	finding, err := MapAlertToOCSF(rawAlert)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if finding.Message != "SSH Brute Force attempt" {
		t.Errorf("Expected message 'SSH Brute Force attempt', got '%s'", finding.Message)
	}

	if finding.Severity != "High" {
		t.Errorf("Expected severity 'High', got '%s'", finding.Severity)
	}

	if finding.SeverityID != 4 {
		t.Errorf("Expected severity ID 4, got %d", finding.SeverityID)
	}

	if len(finding.Observables) != 2 {
		t.Errorf("Expected 2 observables (SrcIP, AgentIP), got %d", len(finding.Observables))
	}

	expectedTime, _ := time.Parse(time.RFC3339, "2023-10-01T12:00:00Z")
	if finding.Time != expectedTime {
		t.Errorf("Expected time %v, got %v", expectedTime, finding.Time)
	}
}

func TestMapAlertToOCSF_InvalidJSON(t *testing.T) {
	_, err := MapAlertToOCSF([]byte(`{invalid json`))
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}
