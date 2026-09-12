package zeek

import (
	"testing"
)

func TestMapZeekConnToOCSF(t *testing.T) {
	rawJSON := []byte(`{"ts":1690000000.123,"uid":"C1234567","id.orig_h":"192.168.1.5","id.orig_p":53421,"id.resp_h":"8.8.8.8","id.resp_p":53,"proto":"udp","service":"dns"}`)
	
	ocsf, err := MapZeekConnToOCSF(rawJSON)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if ocsf.ActivityName != "Network Connection" {
		t.Errorf("expected Network Connection, got %s", ocsf.ActivityName)
	}

	if ocsf.Observables[0].Value != "192.168.1.5" {
		t.Errorf("expected src_ip 192.168.1.5, got %s", ocsf.Observables[0].Value)
	}
}

func TestMapZeekDNSToOCSF(t *testing.T) {
	rawJSON := []byte(`{"ts":1690000000.456,"uid":"C9876543","id.orig_h":"10.0.0.5","id.orig_p":12345,"id.resp_h":"1.1.1.1","id.resp_p":53,"proto":"udp","query":"example.com","qtype":1,"rcode":0}`)
	
	ocsf, err := MapZeekDNSToOCSF(rawJSON)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if ocsf.ActivityName != "DNS Query" {
		t.Errorf("expected DNS Query, got %s", ocsf.ActivityName)
	}

	if ocsf.Observables[0].Value != "example.com" {
		t.Errorf("expected query example.com, got %s", ocsf.Observables[0].Value)
	}
}

func TestMapZeekConn_Malformed(t *testing.T) {
	rawJSON := []byte(`{"ts":1690000000.123,"uid":"C1234567","id.orig_h":12345`) // invalid JSON
	
	_, err := MapZeekConnToOCSF(rawJSON)
	if err == nil {
		t.Fatal("expected error for malformed json, got nil")
	}
}

func TestMapZeekConn_MissingFields(t *testing.T) {
	// Missing proto and IPs
	rawJSON := []byte(`{"ts":1690000000.123,"uid":"C1234567"}`)
	
	ocsf, err := MapZeekConnToOCSF(rawJSON)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if ocsf.Observables[0].Value != "" {
		t.Errorf("expected empty src_ip due to missing field, got %s", ocsf.Observables[0].Value)
	}
}
