package cortex

import (
	"encoding/json"
	"fmt"
	"sae-core/models"
	"time"
)

func MapJobToOCSF(rawEvent []byte) ([]models.OCSFFinding, error) {
	var job JobReport
	if err := json.Unmarshal(rawEvent, &job); err != nil {
		return nil, fmt.Errorf("failed to parse Cortex Job: %w", err)
	}

	if job.ID == "" && job.AnalyzerId == "" {
		return nil, fmt.Errorf("empty or invalid Cortex job report")
	}

	severity := "Info"
	if job.Status == "Failure" {
		severity = "High"
	}

	ocsf := models.OCSFFinding{
		Time:         time.Now(),
		ActivityName: "Threat Analysis Result",
		Severity:     severity,
	}

	if job.Date > 0 {
		ocsf.Time = time.Unix(0, job.Date*int64(time.Millisecond))
	}

	ocsf.Metadata.Product = "SAE Tool (Cortex Engine)"

	ocsf.Message = fmt.Sprintf("Cortex Analysis [%s] completed with status %s", job.AnalyzerName, job.Status)
	if job.AnalyzerName == "" {
		ocsf.Message = fmt.Sprintf("Cortex Analysis [%s] completed with status %s", job.AnalyzerId, job.Status)
	}

	if job.ID != "" {
		ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "JobID", Value: job.ID})
	}
	if job.AnalyzerId != "" {
		ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "AnalyzerID", Value: job.AnalyzerId})
	}

	if job.Report != nil && len(job.Report) > 0 {
		// Just attach a boolean or string indicating report was generated
		ocsf.Observables = append(ocsf.Observables, models.Observable{Type: "HasReport", Value: "true"})
	}

	return []models.OCSFFinding{ocsf}, nil
}

