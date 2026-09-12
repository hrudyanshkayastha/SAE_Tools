package ollama

import (
	"encoding/json"
	"fmt"
	"sae-core/models"
	"time"
)

func ParseAIResponseToOCSF(genResp *GenerateResponse) ([]models.OCSFFinding, error) {
	if genResp == nil || genResp.Response == "" {
		return nil, fmt.Errorf("empty AI response")
	}

	var aiResult AIResult
	// Attempt to extract raw JSON from the LLM response
	// The prompt strictly asks for JSON, but models sometimes wrap it in markdown block.
	responseStr := genResp.Response
	
	err := json.Unmarshal([]byte(responseStr), &aiResult)
	if err != nil {
		// Just store the raw output if it failed to parse as valid JSON
		aiResult.Recommendation = responseStr
		aiResult.Severity = "Info"
		aiResult.RiskScore = 0
	}

	severity := aiResult.Severity
	if severity == "" {
		severity = "Info"
	}

	ocsf := models.OCSFFinding{
		Time:         time.Now(),
		ActivityName: "AI Threat Investigation",
		Severity:     severity,
		Message:      fmt.Sprintf("SAE Local AI Runtime evaluation: %s", aiResult.Recommendation),
	}

	ocsf.Metadata.Product = "SAE Tool (Ollama Runtime)"

	ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "RiskScore", Value: fmt.Sprintf("%d", aiResult.RiskScore)})
	ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "Model", Value: genResp.Model})

	return []models.OCSFFinding{ocsf}, nil
}
