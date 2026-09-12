package thehive

import (
	"encoding/json"
	"fmt"
	"sae-core/models"
	"strings"
	"time"
)

func mapSeverityToOCSF(severity int) string {
	switch severity {
	case 1:
		return "Low"
	case 2:
		return "Medium"
	case 3:
		return "High"
	case 4:
		return "Critical"
	default:
		return "Info"
	}
}

func MapCaseToOCSF(rawEvent []byte) ([]models.OCSFFinding, error) {
	// Attempt to parse as a webhook wrapper first
	var wrapper WebhookEvent
	var c Case

	if err := json.Unmarshal(rawEvent, &wrapper); err == nil && wrapper.Object.ID != "" {
		c = wrapper.Object
	} else {
		// Fallback to flat case JSON
		if err := json.Unmarshal(rawEvent, &c); err != nil {
			return nil, fmt.Errorf("failed to parse TheHive Case: %w", err)
		}
	}

	if c.ID == "" && c.Title == "" {
		return nil, fmt.Errorf("empty or invalid TheHive case")
	}

	ocsf := models.OCSFFinding{
		Time:         time.Now(),
		ActivityName: "Case Management",
		Severity:     mapSeverityToOCSF(c.Severity),
	}

	if c.StartDate > 0 {
		ocsf.Time = time.Unix(0, c.StartDate*int64(time.Millisecond))
	}

	ocsf.Metadata.Product = "SAE Tool (TheHive Engine)"

	desc := c.Description
	if len(desc) > 100 {
		desc = desc[:97] + "..."
	}

	ocsf.Message = fmt.Sprintf("Case #%d: %s - %s", c.Number, c.Title, desc)

	if c.ID != "" {
		ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "CaseID", Value: c.ID})
	}
	if c.Status != "" {
		ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "Status", Value: c.Status})
	}
	if c.Assignee != "" {
		ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "Assignee", Value: c.Assignee})
	}
	if len(c.Tags) > 0 {
		ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "Tags", Value: strings.Join(c.Tags, ",")})
	}

	return []models.OCSFFinding{ocsf}, nil
}
