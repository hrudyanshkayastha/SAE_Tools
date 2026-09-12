# SAE Garak Build Evidence

## 1. Native Build Environment
- **Environment**: Native Windows Host Environment.
- **Engine**: Python 3.x with `garak` v0.16.0 installed.

## 2. Blockers
- **None**: Because Garak is an offline pip module accessible directly via Python invocation, there were no blocked dependencies. It was natively invoked via the `garak` binary from the CLI.

## 3. Results
- The integration executes cleanly.
- Overcame a Windows terminal charset encoding issue (`charmap` codec) by forcing `PYTHONIOENCODING="utf-8"`.
- A real LLM scan report was generated locally at `~/.local/share/garak/garak_runs/`.
