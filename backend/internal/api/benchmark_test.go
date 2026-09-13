package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TEST 6 - API PERFORMANCE
func BenchmarkAPI_Health(b *testing.B) {
	req, _ := http.NewRequest("GET", "/health", nil)
	server := NewServer(nil)
	handler := http.HandlerFunc(server.handleHealth)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkAPI_Dashboard(b *testing.B) {
	req, _ := http.NewRequest("GET", "/dashboard", nil)
	server := NewServer(nil)
	handler := http.HandlerFunc(server.handleDashboard)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}
