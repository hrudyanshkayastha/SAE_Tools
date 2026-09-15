# SAE Ollama Source Audit

## 1. Local Runtime Engine
- **Engine**: Ollama natively installed on the host OS.
- **Version**: 0.17.1
- **Models Loaded**: `llama3.2:3b`, `qwen2.5:7b`, `mistral:latest`, `llama3:latest`, `Foundation-Sec-8B-GGUF`.

## 2. Integration Type
Unlike other tools that required codebase cloning, Ollama is utilized as a precompiled daemon exposing the local REST API on `127.0.0.1:11434`. The Go client communicates directly with this API, requiring no source code modification to Ollama itself.
