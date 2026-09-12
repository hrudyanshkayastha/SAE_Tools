package langgraph

import (
	"fmt"
	"sae-core/models"
	"time"
)

func ParseGraphResultToOCSF(result *GraphOutput) ([]models.OCSFFinding, error) {
	if result == nil {
		return nil, fmt.Errorf("nil result")
	}

	severity := "Info"
	if result.RiskScore >= 80 {
		severity = "High"
	} else if result.RiskScore >= 50 {
		severity = "Medium"
	} else if result.RiskScore >= 30 {
		severity = "Low"
	}

	msg := fmt.Sprintf("SAE Reasoning Engine Decision: %s. Validation: %s.", result.Decision, result.Validation)

	ocsf := models.OCSFFinding{
		Time:         time.Now(),
		ActivityName: "SAE Reasoning Workflow",
		Severity:     severity,
		Message:      msg,
	}

	ocsf.Metadata.Product = "SAE Tool (LangGraph Engine)"

	ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "RiskScore", Value: fmt.Sprintf("%d", result.RiskScore)})
	ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "Decision", Value: result.Decision})

	if result.LLMRecommendation != nil {
		if rec, ok := result.LLMRecommendation["recommendation"].(string); ok {
			ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "LLMRecommendation", Value: rec})
		}
	}

	return []models.OCSFFinding{ocsf}, nil
}
