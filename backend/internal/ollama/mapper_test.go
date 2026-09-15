package ollama

import (
	"strings"
	"testing"
)

func TestParseAIResponseToOCSF_ValidJSON(t *testing.T) {
	resp := &GenerateResponse{
		Model:    "llama3",
		Response: `{"recommendation":"Block IP","severity":"High","risk_score":95}`,
	}
	f, err := ParseAIResponseToOCSF(resp)
	if err != nil || len(f) != 1 { t.Fatalf("expected 1 finding") }
	if f[0].Severity != "High" { t.Errorf("expected High") }
}

func TestParseAIResponseToOCSF_InvalidJSONFallback(t *testing.T) {
	resp := &GenerateResponse{
		Model:    "llama3",
		Response: `I recommend blocking the IP because it's malicious.`,
	}
	f, _ := ParseAIResponseToOCSF(resp)
	if f[0].Severity != "Info" { t.Errorf("expected fallback Info") }
	if !strings.Contains(f[0].Message, "block") { t.Errorf("message missing content") }
}

func TestParseAIResponseToOCSF_EmptyResponse(t *testing.T) {
	_, err := ParseAIResponseToOCSF(nil)
	if err == nil { t.Errorf("expected error on nil response") }
}

func TestParseAIResponseToOCSF_EmptyString(t *testing.T) {
	resp := &GenerateResponse{Response: ""}
	_, err := ParseAIResponseToOCSF(resp)
	if err == nil { t.Errorf("expected error on empty string") }
}

func TestParseAIResponseToOCSF_ExtractsRiskScore(t *testing.T) {
	resp := &GenerateResponse{
		Model:    "llama3",
		Response: `{"recommendation":"Monitor","severity":"Medium","risk_score":50}`,
	}
	f, _ := ParseAIResponseToOCSF(resp)
	found := false
	for _, o := range f[0].Observables {
		if o.Type == "RiskScore" && o.Value == "50" {
			found = true
		}
	}
	if !found { t.Errorf("missing RiskScore observable") }
}

func TestParseAIResponseToOCSF_ExtractsModel(t *testing.T) {
	resp := &GenerateResponse{
		Model:    "llama-sec:latest",
		Response: `{"recommendation":"Monitor"}`,
	}
	f, _ := ParseAIResponseToOCSF(resp)
	found := false
	for _, o := range f[0].Observables {
		if o.Type == "Model" && o.Value == "llama-sec:latest" {
			found = true
		}
	}
	if !found { t.Errorf("missing Model observable") }
}

func TestParseAIResponseToOCSF_ActivityName(t *testing.T) {
	resp := &GenerateResponse{Response: "{}"}
	f, _ := ParseAIResponseToOCSF(resp)
	if f[0].ActivityName != "AI Threat Investigation" { t.Errorf("bad ActivityName") }
}

func TestParseAIResponseToOCSF_ProductMetadata(t *testing.T) {
	resp := &GenerateResponse{Response: "{}"}
	f, _ := ParseAIResponseToOCSF(resp)
	if f[0].Metadata.Product != "SAE Tool (Ollama Runtime)" { t.Errorf("bad Product") }
}

func TestParseAIResponseToOCSF_SeverityCritical(t *testing.T) {
	resp := &GenerateResponse{Response: `{"severity":"Critical"}`}
	f, _ := ParseAIResponseToOCSF(resp)
	if f[0].Severity != "Critical" { t.Errorf("expected Critical") }
}

func TestParseAIResponseToOCSF_SeverityLow(t *testing.T) {
	resp := &GenerateResponse{Response: `{"severity":"Low"}`}
	f, _ := ParseAIResponseToOCSF(resp)
	if f[0].Severity != "Low" { t.Errorf("expected Low") }
}

func TestParseAIResponseToOCSF_MissingSeverity(t *testing.T) {
	resp := &GenerateResponse{Response: `{"recommendation":"test"}`}
	f, _ := ParseAIResponseToOCSF(resp)
	if f[0].Severity != "Info" { t.Errorf("expected Info default") }
}

func TestParseAIResponseToOCSF_DefaultRiskScore(t *testing.T) {
	resp := &GenerateResponse{Response: `{"recommendation":"test"}`}
	f, _ := ParseAIResponseToOCSF(resp)
	for _, o := range f[0].Observables {
		if o.Type == "RiskScore" && o.Value != "0" {
			t.Errorf("expected default risk score 0")
		}
	}
}

func TestParseAIResponseToOCSF_MalformedJSONString(t *testing.T) {
	resp := &GenerateResponse{Response: `{"severity":"High", `}
	f, _ := ParseAIResponseToOCSF(resp)
	if f[0].Severity != "Info" { t.Errorf("should fallback to Info if invalid json") }
}

func TestParseAIResponseToOCSF_TypeMismatchJSON(t *testing.T) {
	resp := &GenerateResponse{Response: `{"severity":10}`}
	f, _ := ParseAIResponseToOCSF(resp)
	if f[0].Severity != "Info" { t.Errorf("type mismatch on severity should fallback to Info") }
}

func TestParseAIResponseToOCSF_TypeMismatchRisk(t *testing.T) {
	resp := &GenerateResponse{Response: `{"risk_score":"high"}`}
	f, _ := ParseAIResponseToOCSF(resp)
	for _, o := range f[0].Observables {
		if o.Type == "RiskScore" && o.Value != "0" {
			t.Errorf("expected default risk score 0 due to type mismatch")
		}
	}
}

func TestParseAIResponseToOCSF_NoObservablesMissing(t *testing.T) {
	resp := &GenerateResponse{Response: `{"recommendation":"test"}`}
	f, _ := ParseAIResponseToOCSF(resp)
	if len(f[0].Observables) != 2 { t.Errorf("expected exactly 2 observables") }
}

func TestParseAIResponseToOCSF_MessagePrefix(t *testing.T) {
	resp := &GenerateResponse{Response: `{"recommendation":"hello"}`}
	f, _ := ParseAIResponseToOCSF(resp)
	if !strings.HasPrefix(f[0].Message, "SAE Local AI Runtime evaluation: ") { t.Errorf("missing standard prefix") }
}

func TestParseAIResponseToOCSF_WhitespaceHandling(t *testing.T) {
	resp := &GenerateResponse{Response: `   {"recommendation":"hello"}   `}
	f, _ := ParseAIResponseToOCSF(resp)
	if !strings.Contains(f[0].Message, "hello") { t.Errorf("failed to parse padded json") }
}
