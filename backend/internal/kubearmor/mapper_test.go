package kubearmor

import (
	"strings"
	"testing"
)

func TestMapAlertToOCSF_ValidRuntime(t *testing.T) {
	raw := `{"Timestamp":16000000,"ClusterName":"default","HostName":"node1","Operation":"File","Resource":"/etc/passwd","Action":"Audit"}`
	ocsf, err := MapAlertToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ocsf.Severity != "Info" {
		t.Errorf("expected Info severity for missing/unknown severity")
	}
	if !strings.Contains(ocsf.Message, "File on /etc/passwd") {
		t.Errorf("expected operation and resource in message")
	}
}

func TestMapAlertToOCSF_PolicyViolation(t *testing.T) {
	raw := `{"PolicyName":"Block-Exec","Operation":"Process","Resource":"/bin/bash","Action":"Block","Severity":"1"}`
	ocsf, err := MapAlertToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error")
	}
	if ocsf.Severity != "Fatal" {
		t.Errorf("expected Fatal severity")
	}
	if !strings.Contains(ocsf.Message, "Block-Exec") {
		t.Errorf("expected policy name in message")
	}
}

func TestMapAlertToOCSF_ProcessEvent(t *testing.T) {
	raw := `{"Operation":"Process","Resource":"/bin/ls","ProcessName":"/bin/ls"}`
	ocsf, err := MapAlertToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error")
	}
	found := false
	for _, obs := range ocsf.Observables {
		if obs.Type == "Process" && obs.Value == "/bin/ls" {
			found = true
		}
	}
	if !found {
		t.Errorf("missing process observable")
	}
}

func TestMapAlertToOCSF_FileEvent(t *testing.T) {
	raw := `{"Operation":"File","Resource":"/etc/shadow"}`
	ocsf, err := MapAlertToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error")
	}
	found := false
	for _, obs := range ocsf.Observables {
		if obs.Type == "File" && obs.Value == "/etc/shadow" {
			found = true
		}
	}
	if !found {
		t.Errorf("missing file observable")
	}
}

func TestMapAlertToOCSF_NetworkEvent(t *testing.T) {
	raw := `{"Operation":"Network","Resource":"TCP"}`
	ocsf, err := MapAlertToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error")
	}
	found := false
	for _, obs := range ocsf.Observables {
		if obs.Type == "File" && obs.Value == "TCP" {
			found = true
		}
	}
	if !found {
		t.Errorf("missing network observable")
	}
}

func TestMapAlertToOCSF_ContainerMetadata(t *testing.T) {
	raw := `{"Operation":"File","ContainerID":"docker-1234","ContainerName":"ubuntu"}`
	ocsf, err := MapAlertToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error")
	}
	found := false
	for _, obs := range ocsf.Observables {
		if obs.Type == "Container" && obs.Value == "docker-1234" {
			found = true
		}
	}
	if !found {
		t.Errorf("missing container observable")
	}
}

func TestMapAlertToOCSF_KubernetesMetadata(t *testing.T) {
	raw := `{"Operation":"File","NamespaceName":"default","PodName":"nginx-pod"}`
	ocsf, err := MapAlertToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error")
	}
	foundNs, foundPod := false, false
	for _, obs := range ocsf.Observables {
		if obs.Type == "Namespace" && obs.Value == "default" {
			foundNs = true
		}
		if obs.Type == "Pod" && obs.Value == "nginx-pod" {
			foundPod = true
		}
	}
	if !foundNs || !foundPod {
		t.Errorf("missing k8s observables")
	}
}

func TestMapAlertToOCSF_MissingOptional(t *testing.T) {
	raw := `{"Operation":"File"}`
	ocsf, err := MapAlertToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error")
	}
	if ocsf.Severity != "Info" {
		t.Errorf("expected Info severity")
	}
}

func TestMapAlertToOCSF_MissingRequired(t *testing.T) {
	raw := `{"Timestamp":16000000}`
	_, err := MapAlertToOCSF([]byte(raw))
	if err == nil {
		t.Fatalf("expected error for missing required fields")
	}
}

func TestMapAlertToOCSF_MalformedJSON(t *testing.T) {
	raw := `{"Operation":"File",`
	_, err := MapAlertToOCSF([]byte(raw))
	if err == nil {
		t.Fatalf("expected error for malformed JSON")
	}
}

func TestMapAlertToOCSF_InvalidFieldTypes(t *testing.T) {
	raw := `{"Operation":123}`
	_, err := MapAlertToOCSF([]byte(raw))
	if err == nil {
		t.Fatalf("expected error for invalid field types")
	}
}

func TestMapAlertToOCSF_InvalidTimestamp(t *testing.T) {
	raw := `{"Timestamp":"invalid","Operation":"File"}`
	_, err := MapAlertToOCSF([]byte(raw))
	if err == nil {
		t.Fatalf("expected error for invalid timestamp type")
	}
}

func TestMapAlertToOCSF_UnknownAction(t *testing.T) {
	raw := `{"Operation":"File","Action":"UnknownAction"}`
	ocsf, err := MapAlertToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error")
	}
	if !strings.Contains(ocsf.Message, "UnknownAction") {
		t.Errorf("expected unknown action to be in message")
	}
}

func TestMapAlertToOCSF_UnknownEventType(t *testing.T) {
	raw := `{"Operation":"WeirdOp","Type":"WeirdType"}`
	ocsf, err := MapAlertToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error")
	}
	if !strings.Contains(ocsf.Message, "WeirdOp") {
		t.Errorf("expected weird operation to be in message")
	}
}
