package suricata

import (
	"strings"
	"testing"
)

func TestMapEVEToOCSF_Alert(t *testing.T) {
	raw := `{"timestamp":"2023-10-12T10:00:00.000Z","event_type":"alert","src_ip":"192.168.1.100","src_port":54321,"dest_ip":"10.0.0.1","dest_port":80,"proto":"TCP","alert":{"action":"allowed","gid":1,"signature_id":2010935,"rev":3,"signature":"ET EXPLOIT Possible CVE-2020-0601","category":"Attempted Administrator Privilege Gain","severity":1}}`
	ocsf, err := MapEVEToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ocsf.ActivityName != "Intrusion Detection" {
		t.Errorf("expected Intrusion Detection, got %s", ocsf.ActivityName)
	}
	if ocsf.Severity != "Critical" {
		t.Errorf("expected Critical severity (1), got %s", ocsf.Severity)
	}
	if !strings.Contains(ocsf.Message, "CVE-2020-0601") {
		t.Errorf("expected signature in message, got %s", ocsf.Message)
	}
}

func TestMapEVEToOCSF_Flow(t *testing.T) {
	raw := `{"timestamp":"2023-10-12T10:00:00.000Z","event_type":"flow","src_ip":"192.168.1.100","src_port":54321,"dest_ip":"10.0.0.1","dest_port":80,"proto":"TCP","flow":{"pkts_toserver":5,"pkts_toclient":5,"bytes_toserver":500,"bytes_toclient":500,"start":"2023-10-12T09:59:00.000Z","end":"2023-10-12T10:00:00.000Z","age":60,"state":"established","reason":"timeout","alerted":false}}`
	ocsf, err := MapEVEToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ocsf.ActivityName != "Network Activity" {
		t.Errorf("expected Network Activity, got %s", ocsf.ActivityName)
	}
}

func TestMapEVEToOCSF_DNS(t *testing.T) {
	raw := `{"timestamp":"2023-10-12T10:00:00.000Z","event_type":"dns","src_ip":"192.168.1.100","src_port":54321,"dest_ip":"8.8.8.8","dest_port":53,"proto":"UDP","dns":{"type":"query","id":1234,"rrname":"example.com","rrtype":"A","rcode":"NOERROR"}}`
	ocsf, err := MapEVEToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ocsf.ActivityName != "DNS Activity" {
		t.Errorf("expected DNS Activity, got %s", ocsf.ActivityName)
	}
	if !strings.Contains(ocsf.Message, "example.com") {
		t.Errorf("expected rrname in message, got %s", ocsf.Message)
	}
}

func TestMapEVEToOCSF_HTTP(t *testing.T) {
	raw := `{"timestamp":"2023-10-12T10:00:00.000Z","event_type":"http","src_ip":"192.168.1.100","src_port":54321,"dest_ip":"10.0.0.1","dest_port":80,"proto":"TCP","http":{"hostname":"example.com","url":"/index.html","http_user_agent":"Mozilla/5.0","http_method":"GET","protocol":"HTTP/1.1","status":200,"length":1024}}`
	ocsf, err := MapEVEToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ocsf.ActivityName != "HTTP Activity" {
		t.Errorf("expected HTTP Activity, got %s", ocsf.ActivityName)
	}
}

func TestMapEVEToOCSF_TLS(t *testing.T) {
	raw := `{"timestamp":"2023-10-12T10:00:00.000Z","event_type":"tls","src_ip":"192.168.1.100","src_port":54321,"dest_ip":"10.0.0.1","dest_port":443,"proto":"TCP","tls":{"subject":"C=US, ST=CA, L=San Francisco, O=Example, CN=example.com","issuerdn":"C=US, O=Let's Encrypt, CN=Let's Encrypt Authority X3","serial":"1234567890","fingerprint":"12:34:56:78:90","sni":"example.com","version":"TLSv1.2"}}`
	ocsf, err := MapEVEToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ocsf.ActivityName != "TLS Activity" {
		t.Errorf("expected TLS Activity, got %s", ocsf.ActivityName)
	}
	if !strings.Contains(ocsf.Message, "example.com") {
		t.Errorf("expected sni in message, got %s", ocsf.Message)
	}
}

func TestMapEVEToOCSF_Malformed(t *testing.T) {
	raw := `{"timestamp":"2023-10-12T10:00:00.000Z","event_type":"alert",`
	_, err := MapEVEToOCSF([]byte(raw))
	if err == nil {
		t.Fatalf("expected error for malformed json")
	}
}

func TestMapEVEToOCSF_MissingFields(t *testing.T) {
	raw := `{"event_type":"alert"}`
	ocsf, err := MapEVEToOCSF([]byte(raw))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ocsf.Severity != "High" {
		t.Errorf("expected fallback severity High, got %s", ocsf.Severity)
	}
}
