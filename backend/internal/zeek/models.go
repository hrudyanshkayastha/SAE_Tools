package zeek

// ZeekConnRecord represents a line from Zeek's conn.log in JSON format.
type ZeekConnRecord struct {
	Ts             float64 `json:"ts"`
	Uid            string  `json:"uid"`
	IdOrigH        string  `json:"id.orig_h"`
	IdOrigP        int     `json:"id.orig_p"`
	IdRespH        string  `json:"id.resp_h"`
	IdRespP        int     `json:"id.resp_p"`
	Proto          string  `json:"proto"`
	Service        string  `json:"service,omitempty"`
	Duration       float64 `json:"duration,omitempty"`
	OrigBytes      int64   `json:"orig_bytes,omitempty"`
	RespBytes      int64   `json:"resp_bytes,omitempty"`
	ConnState      string  `json:"conn_state,omitempty"`
	LocalOrig      bool    `json:"local_orig,omitempty"`
	LocalResp      bool    `json:"local_resp,omitempty"`
	MissedBytes    int64   `json:"missed_bytes,omitempty"`
	History        string  `json:"history,omitempty"`
	OrigPkts       int64   `json:"orig_pkts,omitempty"`
	OrigIpBytes    int64   `json:"orig_ip_bytes,omitempty"`
	RespPkts       int64   `json:"resp_pkts,omitempty"`
	RespIpBytes    int64   `json:"resp_ip_bytes,omitempty"`
	TunnelParents  []string `json:"tunnel_parents,omitempty"`
}

// ZeekDNSRecord represents a line from Zeek's dns.log in JSON format.
type ZeekDNSRecord struct {
	Ts             float64 `json:"ts"`
	Uid            string  `json:"uid"`
	IdOrigH        string  `json:"id.orig_h"`
	IdOrigP        int     `json:"id.orig_p"`
	IdRespH        string  `json:"id.resp_h"`
	IdRespP        int     `json:"id.resp_p"`
	Proto          string  `json:"proto"`
	TransId        int     `json:"trans_id,omitempty"`
	Rtt            float64 `json:"rtt,omitempty"`
	Query          string  `json:"query,omitempty"`
	Qclass         int     `json:"qclass,omitempty"`
	QclassName     string  `json:"qclass_name,omitempty"`
	Qtype          int     `json:"qtype,omitempty"`
	QtypeName      string  `json:"qtype_name,omitempty"`
	Rcode          int     `json:"rcode,omitempty"`
	RcodeName      string  `json:"rcode_name,omitempty"`
	AA             bool    `json:"AA,omitempty"`
	TC             bool    `json:"TC,omitempty"`
	RD             bool    `json:"RD,omitempty"`
	RA             bool    `json:"RA,omitempty"`
	Z              int     `json:"Z,omitempty"`
	Answers        []string `json:"answers,omitempty"`
	TTLs           []float64 `json:"TTLs,omitempty"`
	Rejected       bool    `json:"rejected,omitempty"`
}

// ZeekHTTPRecord represents a line from Zeek's http.log in JSON format.
type ZeekHTTPRecord struct {
	Ts             float64 `json:"ts"`
	Uid            string  `json:"uid"`
	IdOrigH        string  `json:"id.orig_h"`
	IdOrigP        int     `json:"id.orig_p"`
	IdRespH        string  `json:"id.resp_h"`
	IdRespP        int     `json:"id.resp_p"`
	TransDepth     int     `json:"trans_depth"`
	Method         string  `json:"method"`
	Host           string  `json:"host"`
	Uri            string  `json:"uri"`
	Referrer       string  `json:"referrer,omitempty"`
	UserAgent      string  `json:"user_agent,omitempty"`
	ReqBodyLen     int64   `json:"request_body_len,omitempty"`
	RespBodyLen    int64   `json:"response_body_len,omitempty"`
	StatusCode     int     `json:"status_code,omitempty"`
	StatusMsg      string  `json:"status_msg,omitempty"`
}

// GenericZeekRecord to determine the type before full parsing
type GenericZeekRecord map[string]interface{}
