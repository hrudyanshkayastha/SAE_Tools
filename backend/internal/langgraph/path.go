package langgraph

import (
	"os"
	"path/filepath"
)

// GetProjectRoot finds the root directory containing go.mod
func GetProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}

// GetGraphScriptPath returns the absolute path to graph.py
func GetGraphScriptPath() string {
	return filepath.Join(GetProjectRoot(), "internal", "langgraph", "graph.py")
}
