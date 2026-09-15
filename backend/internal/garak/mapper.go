package garak

import (
	"encoding/json"
	"fmt"
	"sae-core/models"
	"time"
)

func ParseGarakEvalToOCSF(rawLine []byte) ([]models.OCSFFinding, error) {
	if len(rawLine) == 0 {
		return nil, fmt.Errorf("empty garak report line")
	}

	// Just in case it's not the eval line
	var generic map[string]interface{}
	if err := json.Unmarshal(rawLine, &generic); err != nil {
		return nil, fmt.Errorf("failed to parse garak line: %w", err)
	}

	entryType, ok := generic["entry_type"].(string)
	if !ok || entryType != "eval" {
		return nil, nil // Ignore non-eval lines
	}

	var eval GarakEval
	if err := json.Unmarshal(rawLine, &eval); err != nil {
		return nil, fmt.Errorf("failed to map garak eval: %w", err)
	}

	var findings []models.OCSFFinding

	for groupName, groupData := range eval.Eval {
		for itemName, itemData := range groupData {
			if itemName == "_summary" {
				continue
			}
			
			details, ok := itemData.(map[string]interface{})
			if !ok {
				continue
			}
			
			detectorName, ok := details["detector_name"].(string)
			if !ok {
				// Might be a probe summary
				continue
			}
			
			detectorDescr, _ := details["detector_descr"].(string)
			absoluteDefcon, _ := details["absolute_defcon"].(float64)
			passed, _ := details["passed"].(float64)
			total, _ := details["total_evaluated"].(float64)

			severity := "Info"
			if passed < total {
				if absoluteDefcon <= 2 {
					severity = "Critical"
				} else if absoluteDefcon == 3 {
					severity = "High"
				} else if absoluteDefcon == 4 {
					severity = "Medium"
				} else {
					severity = "Low"
				}
			}

			ocsf := models.OCSFFinding{
				Time:         time.Now(),
				ActivityName: "AI Vulnerability Scan",
				Severity:     severity,
				Message:      fmt.Sprintf("Garak Detector [%s]: %s (Passed %d/%d)", detectorName, detectorDescr, int(passed), int(total)),
			}

			ocsf.Metadata.Product = "SAE Tool (Garak Scanner)"
			ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "Detector", Value: detectorName})
			ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "Group", Value: groupName})

			findings = append(findings, ocsf)
		}
	}

	return findings, nil
}
