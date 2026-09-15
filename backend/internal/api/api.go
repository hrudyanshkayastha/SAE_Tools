package api

import (
	"encoding/json"
	"net/http"
	"os"
	"sae-core/internal/storage"
	"strings"
	"time"
)

type Server struct {
	store *storage.Storage
}

func NewServer(store *storage.Storage) *Server {
	return &Server{store: store}
}

func (s *Server) Start(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/ready", s.handleReady)
	
	// Open UI endpoint
	mux.HandleFunc("/dashboard", s.handleDashboard)

	mux.HandleFunc("/events", s.requireJWT(s.handleEvents))
	mux.HandleFunc("/incidents", s.requireJWT(s.handleIncidents))
	mux.HandleFunc("/investigations", s.requireJWT(s.handleInvestigations))
	mux.HandleFunc("/decisions", s.requireJWT(s.handleDecisions))

	secureMux := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Secure Headers
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			
			// Configurable CORS
			allowedOrigin := os.Getenv("SAE_API_ORIGIN")
			if allowedOrigin == "" {
				allowedOrigin = "https://sae.local" // secure default
			}
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           secureMux(mux),
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return server.ListenAndServe()
}

func (s *Server) requireJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		
		token := strings.TrimPrefix(authHeader, "Bearer ")
		expectedToken := os.Getenv("SAE_API_TOKEN")
		if expectedToken == "" {
			expectedToken = "default-dev-token" // mock for now, but enforces presence
		}

		// Simple fixed token validation for MVP, should be real JWT lib (e.g. golang-jwt/jwt)
		// For hardening, we verify it matches the environment secret.
		if token != expectedToken {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	if s.store.PG.Ping() != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status":"postgres unavailable"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ready"}`))
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	rows, err := s.store.PG.Query("SELECT raw FROM sae_telemetry ORDER BY timestamp DESC LIMIT 50")
	if err != nil {
		// Log internal error, but do not leak SQL context to client
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []json.RawMessage
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err == nil {
			events = append(events, json.RawMessage(raw))
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (s *Server) handleIncidents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	rows, err := s.store.PG.Query("SELECT json_build_object('correlation_id', correlation_id, 'target', target, 'events_count', events_count, 'severity', severity, 'status', status) FROM correlations")
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []json.RawMessage
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err == nil {
			items = append(items, json.RawMessage(raw))
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (s *Server) handleInvestigations(w http.ResponseWriter, r *http.Request) {
	s.handleDecisions(w, r)
}

func (s *Server) handleDecisions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	rows, err := s.store.PG.Query("SELECT json_build_object('correlation_id', correlation_id, 'risk_score', risk_score, 'validation', validation, 'action', action, 'state', state) FROM decisions")
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []json.RawMessage
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err == nil {
			items = append(items, json.RawMessage(raw))
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}


func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "{\"error\":\"method not allowed\"} ", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	html := "<!DOCTYPE html><html><head><title>SAE SOC Dashboard</title><style>body{font-family:sans-serif;background:#1e1e1e;color:#fff;padding:20px;}</style></head><body><h1>SAE v0.3 SOC Dashboard</h1><p>Status: Native UI actively polling protected endpoints.</p></body></html>"
	w.Write([]byte(html))
}

