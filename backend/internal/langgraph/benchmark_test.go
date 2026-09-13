package langgraph

import (
	"encoding/json"
	"testing"
)

func BenchmarkLangGraph_Inference(b *testing.B) {
	// Measures latency of the real LangGraph/Ollama call
	eventData, _ := json.Marshal(map[string]string{
		"event_source": "Wazuh",
		"event_type":   "Port Scan",
		"description":  "Suspicious port scanning behavior from external IP",
		"severity":     "Medium",
		"source_ip":    "10.0.0.5",
	})
	
	// Benchmark is very slow for LLMs, so we limit b.N to just what go test decides, 
	// typically we just want the time/op.
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ExecuteReasoningGraph(eventData, "E:\\New folder\\SAE_Tools\\SAE\\backend\\internal\\langgraph")
		if err != nil {
			b.Fatalf("LangGraph execution failed: %v", err)
		}
	}
}
