package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPI_Dashboard(t *testing.T) {
	req, _ := http.NewRequest("GET", "/dashboard", nil)
	rr := httptest.NewRecorder()
	server := NewServer(nil)
	
	handler := http.HandlerFunc(server.handleDashboard)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !strings.Contains(rr.Body.String(), "SAE SOC Dashboard") {
		t.Errorf("Dashboard HTML did not contain title, got: %v", rr.Body.String())
	}
}

func TestAPI_AuthRejection(t *testing.T) {
	server := NewServer(nil)
	handler := server.requireJWT(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/events", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler should reject without JWT, got %v", status)
	}
}

func TestAPI_AuthSuccess(t *testing.T) {
	server := NewServer(nil)
	handler := server.requireJWT(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/events", nil)
	req.Header.Set("Authorization", "Bearer default-dev-token")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler should accept valid default JWT token, got %v", status)
	}
}
