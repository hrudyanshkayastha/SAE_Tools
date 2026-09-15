package verification

import (
	"context"
	"time"
)

const (
	StateVerified    = "VERIFIED"
	StateNotVerified = "NOT_VERIFIED"
	StateUnknown     = "UNKNOWN"
	StateError       = "ERROR"
)

type Result struct {
	Action    string
	Target    string
	Expected  string
	Observed  string
	State     string
	Evidence  string
	Timestamp time.Time
}

type TargetVerifier interface {
	Verify(ctx context.Context, action string, target string) Result
}
