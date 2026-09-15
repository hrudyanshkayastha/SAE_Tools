package scoutsuite

type ScoutSuiteReport struct {
	ProviderCode string                          `json:"provider_code"`
	ProviderName string                          `json:"provider_name"`
	AccountID    string                          `json:"account_id"`
	Services     map[string]ScoutSuiteService    `json:"services"`
}

type ScoutSuiteService struct {
	Findings map[string]ScoutSuiteFinding `json:"findings"`
}

type ScoutSuiteFinding struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Rationale    string   `json:"rationale"`
	Remediation  string   `json:"remediation"`
	Level        string   `json:"level"`
	Items        []string `json:"items"`
}
