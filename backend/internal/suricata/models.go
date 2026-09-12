package suricata

import "time"

// SuricataEVE represents the unified JSON structure of Suricata's eve.json
type SuricataEVE struct {
	Timestamp   time.Time `json:"timestamp"`
	EventType   string    `json:"event_type"` // alert, flow, dns, http, tls, etc.
	SrcIP       string    `json:"src_ip"`
	SrcPort     int       `json:"src_port"`
	DestIP      string    `json:"dest_ip"`
	DestPort    int       `json:"dest_port"`
	Proto       string    `json:"proto"`
	AppProto    string    `json:"app_proto"`
	InIface     string    `json:"in_iface"`
	Vlan        []int     `json:"vlan"`
	FlowID      int64     `json:"flow_id"`

	// Alert specifics
	Alert *AlertEvent `json:"alert,omitempty"`

	// Flow specifics
	Flow *FlowEvent `json:"flow,omitempty"`

	// Protocol specifics
	DNS  *DNSEvent  `json:"dns,omitempty"`
	HTTP *HTTPEvent `json:"http,omitempty"`
	TLS  *TLSEvent  `json:"tls,omitempty"`
}

type AlertEvent struct {
	Action      string `json:"action"`
	Gid         int    `json:"gid"`
	SignatureID int    `json:"signature_id"`
	Rev         int    `json:"rev"`
	Signature   string `json:"signature"`
	Category    string `json:"category"`
	Severity    int    `json:"severity"`
	Source      *struct {
		IP   string `json:"ip"`
		Port int    `json:"port"`
	} `json:"source,omitempty"`
	Target      *struct {
		IP   string `json:"ip"`
		Port int    `json:"port"`
	} `json:"target,omitempty"`
}

type FlowEvent struct {
	PktsToserver  int64  `json:"pkts_toserver"`
	PktsToclient  int64  `json:"pkts_toclient"`
	BytesToserver int64  `json:"bytes_toserver"`
	BytesToclient int64  `json:"bytes_toclient"`
	Start         string `json:"start"`
	End           string `json:"end"`
	Age           int    `json:"age"`
	State         string `json:"state"`
	Reason        string `json:"reason"`
	Alerted       bool   `json:"alerted"`
}

type DNSEvent struct {
	Type   string `json:"type"` // query or answer
	ID     int    `json:"id"`
	RRName string `json:"rrname"`
	RRType string `json:"rrtype"`
	RCode  string `json:"rcode"`
}

type HTTPEvent struct {
	Hostname      string `json:"hostname"`
	URL           string `json:"url"`
	HTTPUserAgent string `json:"http_user_agent"`
	HTTPMethod    string `json:"http_method"`
	Protocol      string `json:"protocol"`
	Status        int    `json:"status"`
	Length        int    `json:"length"`
}

type TLSEvent struct {
	Subject string `json:"subject"`
	IssuerDN string `json:"issuerdn"`
	Serial   string `json:"serial"`
	Fingerprint string `json:"fingerprint"`
	Sni      string `json:"sni"`
	Version  string `json:"version"`
}
