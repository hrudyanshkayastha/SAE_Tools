package zeek

import (
	"encoding/json"
	"fmt"
	"sae-core/models"
	"time"
)

// MapZeekConnToOCSF maps a Zeek connection record to an OCSF Network Activity finding.
func MapZeekConnToOCSF(rawJSON []byte) (*models.OCSFFinding, error) {
	var conn ZeekConnRecord
	if err := json.Unmarshal(rawJSON, &conn); err != nil {
		return nil, fmt.Errorf("failed to parse Zeek conn record: %w", err)
	}

	finding := &models.OCSFFinding{
		ActivityID:   3, // Network Activity
		ActivityName: "Network Connection",
		Time:         time.Unix(int64(conn.Ts), int64((conn.Ts-float64(int64(conn.Ts)))*1e9)),
		SeverityID:   1, // Informational by default for standard connections
		Severity:     "Info",
		Status:       "New",
		Message:      fmt.Sprintf("Network connection from %s:%d to %s:%d over %s", conn.IdOrigH, conn.IdOrigP, conn.IdRespH, conn.IdRespP, conn.Proto),
		RawEvent:     string(rawJSON),
	}
	finding.Metadata.Product = "SAE Tool (Zeek Engine)"
	finding.Metadata.Vendor = "Corelight/Zeek"
	finding.Metadata.Version = "9.1.0-dev.93"
	finding.Metadata.LogFormat = "JSON"

	finding.Observables = append(finding.Observables, models.Observable{Name: "src_ip", Type: "IP", Value: conn.IdOrigH})
	finding.Observables = append(finding.Observables, models.Observable{Name: "dst_ip", Type: "IP", Value: conn.IdRespH})

	return finding, nil
}

// MapZeekDNSToOCSF maps a Zeek DNS record to an OCSF finding.
func MapZeekDNSToOCSF(rawJSON []byte) (*models.OCSFFinding, error) {
	var dns ZeekDNSRecord
	if err := json.Unmarshal(rawJSON, &dns); err != nil {
		return nil, fmt.Errorf("failed to parse Zeek dns record: %w", err)
	}

	finding := &models.OCSFFinding{
		ActivityID:   4, // DNS Activity
		ActivityName: "DNS Query",
		Time:         time.Unix(int64(dns.Ts), int64((dns.Ts-float64(int64(dns.Ts)))*1e9)),
		SeverityID:   1,
		Severity:     "Info",
		Status:       "New",
		Message:      fmt.Sprintf("DNS Query for %s from %s", dns.Query, dns.IdOrigH),
		RawEvent:     string(rawJSON),
	}
	finding.Metadata.Product = "SAE Tool (Zeek Engine)"
	finding.Metadata.Vendor = "Corelight/Zeek"
	finding.Metadata.Version = "9.1.0-dev.93"
	finding.Metadata.LogFormat = "JSON"

	finding.Observables = append(finding.Observables, models.Observable{Name: "query", Type: "Domain", Value: dns.Query})
	finding.Observables = append(finding.Observables, models.Observable{Name: "src_ip", Type: "IP", Value: dns.IdOrigH})

	return finding, nil
}
