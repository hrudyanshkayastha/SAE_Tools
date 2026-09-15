package verification

import (
	"context"
	"errors"
	"testing"
)

func TestLocalVerifier_Success(t *testing.T) {
	v := NewLocalTestVerifier()
	ctx := context.Background()

	// Before: NOT_BLOCKED
	res1 := v.Verify(ctx, "BLOCK_IP", "10.0.0.1")
	if res1.State != StateNotVerified || res1.Observed != "NOT_BLOCKED" {
		t.Fatalf("Expected NOT_VERIFIED and NOT_BLOCKED, got %s and %s", res1.State, res1.Observed)
	}

	// Action: BLOCK_IP
	v.BlockIP("10.0.0.1")

	// After: BLOCKED
	res2 := v.Verify(ctx, "BLOCK_IP", "10.0.0.1")
	if res2.State != StateVerified || res2.Observed != "BLOCKED" {
		t.Fatalf("Expected VERIFIED and BLOCKED, got %s and %s", res2.State, res2.Observed)
	}
}

func TestLocalVerifier_FalseSuccess(t *testing.T) {
	v := NewLocalTestVerifier()
	ctx := context.Background()

	// Shuffle returns FINISHED, but target remains NOT_BLOCKED
	res := v.Verify(ctx, "BLOCK_IP", "10.0.0.2")
	
	if res.State != StateNotVerified {
		t.Errorf("Expected FALSE SUCCESS to yield NOT_VERIFIED, got %s", res.State)
	}
}

func TestLocalVerifier_Unknown(t *testing.T) {
	v := NewLocalTestVerifier()
	ctx := context.Background()

	v.SetUnknown(true)
	res := v.Verify(ctx, "BLOCK_IP", "10.0.0.3")

	if res.State != StateUnknown {
		t.Errorf("Expected UNKNOWN, got %s", res.State)
	}
}

func TestLocalVerifier_Error(t *testing.T) {
	v := NewLocalTestVerifier()
	ctx := context.Background()

	v.SetError(errors.New("connection reset by peer"))
	res := v.Verify(ctx, "BLOCK_IP", "10.0.0.4")

	if res.State != StateError {
		t.Errorf("Expected ERROR, got %s", res.State)
	}
}
func TestLocalVerifier_WrongTarget(t *testing.T) {
	v := NewLocalTestVerifier()
	ctx := context.Background()

	// Action hits target A
	v.BlockIP("10.0.0.1")

	// Verify queries target B
	res := v.Verify(ctx, "BLOCK_IP", "10.0.0.5")

	if res.State == StateVerified {
		t.Errorf("Expected NOT_VERIFIED for wrong target, got %s", res.State)
	}
}
