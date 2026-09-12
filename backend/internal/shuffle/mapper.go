package shuffle

import (
	"encoding/json"
	"fmt"
	"sae-core/models"
	"time"
)

func mapStatusToSeverity(status string) string {
	switch status {
	case "SUCCESS":
		return "Info"
	case "FAILURE", "FAILED", "ERROR":
		return "High"
	case "ABORTED":
		return "Medium"
	default:
		return "Info"
	}
}

func MapResultToOCSF(rawEvent []byte) ([]models.OCSFFinding, error) {
	var result ActionResult
	if err := json.Unmarshal(rawEvent, &result); err != nil {
		return nil, fmt.Errorf("failed to parse Shuffle action result: %w", err)
	}

	if result.ExecutionId == "" && result.Action.AppName == "" {
		return nil, fmt.Errorf("empty or invalid Shuffle action result")
	}

	ocsf := models.OCSFFinding{
		Time:         time.Now(),
		ActivityName: "Response Execution",
		Severity:     mapStatusToSeverity(result.Status),
	}

	ocsf.Metadata.Product = "SAE Tool (Shuffle Engine)"

	resText := result.Result
	if len(resText) > 200 {
		resText = resText[:197] + "..."
	}

	ocsf.Message = fmt.Sprintf("Shuffle Action [%s: %s (%s)] completed with status: %s | Result: %s",
		result.Action.AppName, result.Action.Name, result.Action.Label, result.Status, resText)

	if result.ExecutionId != "" {
		ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "ExecutionID", Value: result.ExecutionId})
	}
	if result.Action.AppName != "" {
		ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "AppName", Value: result.Action.AppName})
	}
	if result.Action.Name != "" {
		ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "ActionName", Value: result.Action.Name})
	}
	if result.Status != "" {
		ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "Status", Value: result.Status})
	}

	return []models.OCSFFinding{ocsf}, nil
}
