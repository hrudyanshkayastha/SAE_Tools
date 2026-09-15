package verification

import (
	"context"
	"sync"
	"time"
)

type LocalTestVerifier struct {
	mu         sync.RWMutex
	BlockedIPs map[string]bool
	FailError  error
	IsUnknown  bool
}

func NewLocalTestVerifier() *LocalTestVerifier {
	return &LocalTestVerifier{
		BlockedIPs: make(map[string]bool),
	}
}

func (v *LocalTestVerifier) BlockIP(ip string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.BlockedIPs[ip] = true
}

func (v *LocalTestVerifier) ResetIP(ip string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.BlockedIPs[ip] = false
}

func (v *LocalTestVerifier) SetError(err error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.FailError = err
}

func (v *LocalTestVerifier) SetUnknown(unknown bool) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.IsUnknown = unknown
}

func (v *LocalTestVerifier) Verify(ctx context.Context, action string, target string) Result {
	v.mu.RLock()
	defer v.mu.RUnlock()

	res := Result{
		Action:    action,
		Target:    target,
		Timestamp: time.Now(),
	}

	if v.FailError != nil {
		res.State = StateError
		res.Expected = "KNOWN_STATE"
		res.Observed = "ERROR"
		res.Evidence = v.FailError.Error()
		return res
	}

	if v.IsUnknown {
		res.State = StateUnknown
		res.Expected = "KNOWN_STATE"
		res.Observed = "UNREACHABLE_OR_UNKNOWN"
		res.Evidence = "Target query unavailable"
		return res
	}

	if action == "BLOCK_IP" {
		res.Expected = "BLOCKED"
		if v.BlockedIPs[target] {
			res.Observed = "BLOCKED"
			res.State = StateVerified
			res.Evidence = "IP found in local deterministic blocklist"
		} else {
			res.Observed = "NOT_BLOCKED"
			res.State = StateNotVerified
			res.Evidence = "IP missing from local deterministic blocklist"
		}
		return res
	}

	res.State = StateUnknown
	res.Expected = "UNKNOWN"
	res.Observed = "UNKNOWN"
	res.Evidence = "Unsupported action for LocalTestVerifier"
	return res
}
