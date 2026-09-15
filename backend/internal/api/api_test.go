package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_Health(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	server := &Server{}
	
	handler := http.HandlerFunc(server.handleHealth)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":"ok"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestAPI_ServerInit(t *testing.T) {
	server := NewServer(nil)
	if server == nil {
		t.Errorf("Server should not be nil")
	}
}

func TestAPI_Routes(t *testing.T) {
	server := NewServer(nil)
	if server.store != nil {
		t.Errorf("Store should be nil")
	}
}
