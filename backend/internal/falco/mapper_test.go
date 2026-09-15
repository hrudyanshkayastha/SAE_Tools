package falco

import (
	"strings"
	"testing"
)

func TestMapFalcoToOCSF_ValidRuntime(t *testing.T) {
	raw := `{"output":"Notice shell spawned","priority":"Notice","rule":"Terminal shell","time":"2020-05-01T22:02:40.457813264Z","output_fields":{"user.name":"root"}}`
	ocsf, err := MapFalcoToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ocsf.Severity != "Low" {
		t.Errorf("expected Low severity, got %s", ocsf.Severity)
	}
	if !strings.Contains(ocsf.Message, "Terminal shell") {
		t.Errorf("expected rule in message")
	}
}

func TestMapFalcoToOCSF_HighPriority(t *testing.T) {
	raw := `{"output":"Critical vulnerability exploited","priority":"Critical","rule":"Exploit Attempt","time":"2020-05-01T22:02:40Z","output_fields":{}}`
	ocsf, err := MapFalcoToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error")
	}
	if ocsf.Severity != "Critical" {
		t.Errorf("expected Critical severity")
	}
}

func TestMapFalcoToOCSF_ContainerEvent(t *testing.T) {
	raw := `{"output":"Container escape","priority":"Emergency","rule":"Escape","time":"2020-05-01T22:02:40Z","output_fields":{"container.id":"e8573133241b"}}`
	ocsf, err := MapFalcoToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error")
	}
	found := false
	for _, obs := range ocsf.Observables {
		if obs.Type == "Container" && obs.Value == "e8573133241b" {
			found = true
		}
	}
	if !found {
		t.Errorf("missing container observable")
	}
}

func TestMapFalcoToOCSF_ProcessEvent(t *testing.T) {
	raw := `{"output":"Proc spawn","priority":"Warning","rule":"Spawn","time":"2020-05-01T22:02:40Z","output_fields":{"proc.cmdline":"bash -c 'malicious'"}}`
	ocsf, err := MapFalcoToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error")
	}
	found := false
	for _, obs := range ocsf.Observables {
		if obs.Type == "Process" && obs.Value == "bash -c 'malicious'" {
			found = true
		}
	}
	if !found {
		t.Errorf("missing process observable")
	}
}

func TestMapFalcoToOCSF_FileEvent(t *testing.T) {
	raw := `{"output":"File read","priority":"Warning","rule":"Read Sensitive File","time":"2020-05-01T22:02:40Z","output_fields":{"fd.name":"/etc/shadow"}}`
	ocsf, err := MapFalcoToOCSF([]byte(raw))
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

func TestMapFalcoToOCSF_NetworkEvent(t *testing.T) {
	raw := `{"output":"Network connection","priority":"Notice","rule":"Unexpected network connection","time":"2020-05-01T22:02:40Z","output_fields":{"fd.name":"192.168.1.1:80->10.0.0.5:43212"}}`
	ocsf, err := MapFalcoToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error")
	}
	found := false
	for _, obs := range ocsf.Observables {
		if obs.Type == "File" && obs.Value == "192.168.1.1:80->10.0.0.5:43212" {
			// In Falco, network connections are often represented as fd.name
			found = true
		}
	}
	if !found {
		t.Errorf("missing network fd observable")
	}
}

func TestMapFalcoToOCSF_MissingOptional(t *testing.T) {
	raw := `{"output":"Generic output","priority":"Informational","rule":"Generic rule","time":"2020-05-01T22:02:40Z"}`
	ocsf, err := MapFalcoToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error")
	}
	if ocsf.Severity != "Info" {
		t.Errorf("expected Info severity")
	}
}

func TestMapFalcoToOCSF_MissingRequired(t *testing.T) {
	raw := `{"priority":"Warning","time":"2020-05-01T22:02:40Z"}` // missing rule and output
	_, err := MapFalcoToOCSF([]byte(raw))
	if err == nil {
		t.Fatalf("expected error for missing required fields")
	}
}

func TestMapFalcoToOCSF_MalformedJSON(t *testing.T) {
	raw := `{"output":"Generic output",` // syntax error
	_, err := MapFalcoToOCSF([]byte(raw))
	if err == nil {
		t.Fatalf("expected error for malformed JSON")
	}
}

func TestMapFalcoToOCSF_InvalidFieldTypes(t *testing.T) {
	raw := `{"output":123,"priority":true,"rule":null,"time":"2020-05-01T22:02:40Z"}`
	_, err := MapFalcoToOCSF([]byte(raw))
	if err == nil {
		t.Fatalf("expected error for invalid field types")
	}
}

func TestMapFalcoToOCSF_InvalidTimestamp(t *testing.T) {
	raw := `{"output":"Out","priority":"Warning","rule":"Rule","time":"invalid-time"}`
	_, err := MapFalcoToOCSF([]byte(raw))
	if err == nil {
		t.Fatalf("expected error for invalid time format")
	}
}

func TestMapFalcoToOCSF_UnknownPriority(t *testing.T) {
	raw := `{"output":"Out","priority":"WeirdPriority","rule":"Rule","time":"2020-05-01T22:02:40Z"}`
	ocsf, err := MapFalcoToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error")
	}
	if ocsf.Severity != "Unknown" {
		t.Errorf("expected Unknown severity")
	}
}
