package wazuh

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_AuthenticateAndHealth(t *testing.T) {
	// Mock Wazuh API Server
	mux := http.NewServeMux()
	
	mux.HandleFunc("/security/user/authenticate", func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "wazuh" || pass != "wazuh" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data": {"token": "fake-jwt-token"}}`))
	})

	mux.HandleFunc("/manager/status", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer fake-jwt-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": {"status": "ok"}}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(server.URL, "wazuh", "wazuh", true)
	ctx := context.Background()

	err := client.Authenticate(ctx)
	if err != nil {
		t.Fatalf("Expected no error on auth, got %v", err)
	}

	if client.Token != "fake-jwt-token" {
		t.Errorf("Expected token 'fake-jwt-token', got '%s'", client.Token)
	}

	err = client.HealthCheck(ctx)
	if err != nil {
		t.Fatalf("Expected no error on health check, got %v", err)
	}
}

func TestClient_AuthFailure(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/security/user/authenticate", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(server.URL, "bad", "credentials", true)
	err := client.Authenticate(context.Background())
	if err == nil {
		t.Fatal("Expected error on invalid credentials, got nil")
	}
}
