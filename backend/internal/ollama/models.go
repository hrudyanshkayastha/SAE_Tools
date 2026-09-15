package ollama

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

type AIResult struct {
	Recommendation string `json:"recommendation"`
	Severity       string `json:"severity"`
	RiskScore      int    `json:"risk_score"`
}
