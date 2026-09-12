# SAE Ollama License Audit

## 1. Upstream Identity
- **Project**: Ollama (ollama/ollama)
- **License**: MIT License

## 2. SAE Distribution Posture
- **Integration Boundary**: Ollama is operated as a separate background daemon API. The SAE backend written in Go utilizes HTTP POST requests to query the models.
- **Legal Status**: Because Ollama is MIT-licensed, connecting to it or embedding its binaries carries almost zero legal risk compared to AGPL/GPL components. It seamlessly complies with the SAE modular architecture guidelines.
