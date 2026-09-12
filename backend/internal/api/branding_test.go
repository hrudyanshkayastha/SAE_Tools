package api

import (
	"os"
	"strings"
	"testing"
)

func TestFrontend_SAEBranding(t *testing.T) {
	content, err := os.ReadFile("../../../frontend/src/App.tsx")
	if err != nil {
		t.Skip("Frontend not available in this test environment")
	}

	appCode := string(content)
	if !strings.Contains(appCode, "SAE") {
		t.Errorf("Frontend App.tsx missing SAE branding")
	}
	
	if !strings.Contains(appCode, "SAE Security Fabric") {
		t.Errorf("Frontend App.tsx missing SAE Security Fabric branding")
	}
}

func TestFrontend_ArchitectureSeparation(t *testing.T) {
	content, err := os.ReadFile("../../../frontend/src/pages/Architecture.tsx")
	if err != nil {
		t.Skip("Frontend not available in this test environment")
	}

	archCode := string(content)
	
	expectedModules := []string{
		"Wazuh",
		"Zeek",
		"Falco",
		"LangGraph",
	}

	for _, mod := range expectedModules {
		if !strings.Contains(archCode, mod) {
			t.Errorf("Frontend Architecture missing tool: %s", mod)
		}
	}

	if !strings.Contains(archCode, "FULLY VERIFIED") {
		t.Errorf("Architecture should display tool verification statuses")
	}
}

func TestDocs_ThirdPartyAttribution(t *testing.T) {
	content, err := os.ReadFile("../../../THIRD_PARTY_NOTICES.md")
	if err != nil {
		t.Skip("Notices file not available")
	}

	notices := string(content)
	if !strings.Contains(notices, "GPLv2") || !strings.Contains(notices, "Apache 2.0") {
		t.Errorf("THIRD_PARTY_NOTICES.md missing correct license attributions")
	}
}
