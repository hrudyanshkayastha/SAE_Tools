package models

import "time"

// OCSFFinding represents a standardized Security Finding event in the SAE platform.
// Inspired by the Open Cybersecurity Schema Framework (OCSF).
type OCSFFinding struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	ActivityID    int       `json:"activity_id"`   // e.g., 1 (Create), 2 (Update)
	ActivityName  string    `json:"activity_name"` // e.g., "Alert"
	Time          time.Time `json:"time"`
	SeverityID    int       `json:"severity_id"`   // 1=Info, 2=Low, 3=Medium, 4=High, 5=Critical, 6=Fatal
	Severity      string    `json:"severity"`      // Info, Low, Medium, High, Critical, Fatal
	Status        string    `json:"status"`        // New, In_Progress, Resolved
	Message       string    `json:"message"`
	
	Metadata struct {
		Product   string `json:"product"`
		Vendor    string `json:"vendor"`
		Version   string `json:"version"`
		LogFormat string `json:"log_format"`
	} `json:"metadata"`

	Observables []Observable `json:"observables"`

	// Original raw payload from the source tool
	RawEvent string `json:"raw_event"`
}

type Observable struct {
	Name  string `json:"name"`
	Type  string `json:"type"`  // IP, Hash, Domain, User
	Value string `json:"value"` // 192.168.1.1, admin, etc.
}
