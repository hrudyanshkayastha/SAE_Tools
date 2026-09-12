package falco

import "time"

// FalcoEvent represents the structured JSON output from Falco.
type FalcoEvent struct {
	Output       string                 `json:"output"`
	Priority     string                 `json:"priority"`
	Rule         string                 `json:"rule"`
	Time         time.Time              `json:"time"`
	OutputFields map[string]interface{} `json:"output_fields"`
	Hostname     string                 `json:"hostname,omitempty"`
	Source       string                 `json:"source,omitempty"`
}
