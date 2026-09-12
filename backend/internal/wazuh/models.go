package wazuh

import "time"

// Alert represents the standard Wazuh alerts.json structure
type Alert struct {
	Timestamp   time.Time `json:"timestamp"`
	Rule        Rule      `json:"rule"`
	Agent       Agent     `json:"agent"`
	Manager     Manager   `json:"manager"`
	ID          string    `json:"id"`
	Location    string    `json:"location"`
	Decoder     Decoder   `json:"decoder"`
	FullLog     string    `json:"full_log"`
	Data        Data      `json:"data"`
}

type Rule struct {
	Level       int      `json:"level"`
	Description string   `json:"description"`
	ID          string   `json:"id"`
	FiredTimes  int      `json:"firedtimes"`
	Mail        bool     `json:"mail"`
	Groups      []string `json:"groups"`
	Mitre       Mitre    `json:"mitre"`
}

type Mitre struct {
	ID        []string `json:"id"`
	Tactic    []string `json:"tactic"`
	Technique []string `json:"technique"`
}

type Agent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	IP   string `json:"ip"`
}

type Manager struct {
	Name string `json:"name"`
}

type Decoder struct {
	Name string `json:"name"`
}

type Data struct {
	SrcIP   string `json:"srcip,omitempty"`
	DstIP   string `json:"dstip,omitempty"`
	SrcPort string `json:"srcport,omitempty"`
	DstPort string `json:"dstport,omitempty"`
	Action  string `json:"action,omitempty"`
}

// APIResponse represents a generic Wazuh API response
type APIResponse struct {
	Data   interface{} `json:"data"`
	Error  int         `json:"error"`
	Message string     `json:"message"`
}
