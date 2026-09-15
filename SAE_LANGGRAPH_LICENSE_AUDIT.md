# SAE LangGraph License Audit

## 1. Upstream Identity
- **Project**: LangGraph (langchain-ai/langgraph)
- **License**: MIT License

## 2. SAE Distribution Posture
- **Integration Boundary**: LangGraph is consumed as an unmodified Python library. The SAE application interfaces with it by dynamically executing a `.py` script via Go's `os/exec`.
- **Legal Status**: Because LangGraph is MIT-licensed, importing it natively carries zero legal risk compared to AGPL/GPL components. It poses no viral contamination risk to the Go backend logic.
