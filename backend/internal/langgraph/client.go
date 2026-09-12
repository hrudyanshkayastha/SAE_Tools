package langgraph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
)

type GraphOutput struct {
	Status           string                 `json:"status"`
	Error            string                 `json:"error,omitempty"`
	LLMRecommendation map[string]interface{} `json:"llm_recommendation,omitempty"`
	Validation       string                 `json:"validation,omitempty"`
	RiskScore        int                    `json:"risk_score,omitempty"`
	Decision         string                 `json:"decision,omitempty"`
}

func ExecuteReasoningGraph(eventData []byte, scriptDir string) (*GraphOutput, error) {
	if len(eventData) == 0 {
		return nil, fmt.Errorf("empty event")
	}

	scriptPath := filepath.Join(scriptDir, "graph.py")
	
	cmd := exec.Command("python", scriptPath)
	cmd.Stdin = bytes.NewReader(eventData)
	
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil && out.Len() == 0 {
		return nil, fmt.Errorf("execution failed: %v, stderr: %s", err, stderr.String())
	}
	
	var result GraphOutput
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal graph output: %w", err)
	}
	
	if result.Error != "" {
		return &result, fmt.Errorf("graph returned error: %s", result.Error)
	}
	
	return &result, nil
}
