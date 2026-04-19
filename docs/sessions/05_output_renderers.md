# Session 5: Output Renderers

## Goal
Build the formatting engines for the different ways a user might want to consume the configuration (e.g., .env, console export).

## Tooling & Library Versions
- **Go Version:** `1.26.x`
- **Libraries:** Standard library `fmt`, `strings`, `text/template`.

## Setup Instructions
Create the following structure:
```text
internal/
  render/
    render.go
    env_file.go
    shell_cmd.go
    render_test.go
```

## Tasks

### 1. Define Renderer API
In `internal/render/render.go`:

```go
package render

import "github.com/yourusername/copilot-cli-env/internal/config"

type OutputMode string

const (
	ModeEnvFile OutputMode = "env"
	ModeConsole OutputMode = "console" // direct dump
	ModeCommand OutputMode = "command" // shell-specific commands
)

type Renderer interface {
	Render(env config.CopilotEnv) (string, error)
}
```

### 2. Implement `.env` File Renderer
Format the `CopilotEnv` struct into a standard `.env` format:
```text
COPILOT_PROVIDER_BASE_URL=https://...
COPILOT_MODEL=gpt-4o
COPILOT_PROVIDER_TYPE=openai
COPILOT_PROVIDER_API_KEY=sk-123...
```

### 3. Implement Shell Command Renderer
This needs context about the detected shell (from Session 2).
- **Bash/Zsh:**
  ```bash
  export COPILOT_PROVIDER_BASE_URL="https://..."
  export COPILOT_MODEL="gpt-4o"
  ```
- **PowerShell:**
  ```powershell
  $env:COPILOT_PROVIDER_BASE_URL="https://..."
  $env:COPILOT_MODEL="gpt-4o"
  ```
- **CMD:**
  ```cmd
  set COPILOT_PROVIDER_BASE_URL=https://...
  set COPILOT_MODEL=gpt-4o
  ```

### 4. Create Masking Utility
Create a helper function `MaskSecret(secret string) string` that returns `sk-...****` to be used in UI previews, ensuring the raw string is only used in the final render.

## Unit Testing
Write table-driven unit tests in `render_test.go` to achieve near 100% test coverage.
- Pass a mock `CopilotEnv` to each renderer and assert the exact string output matches the expected Bash/Zsh, PowerShell, or CMD format.
- Verify shell quoting is handled correctly (e.g., escaping quotes if a model name or URL somehow contains them).
- Test that the `MaskSecret` utility correctly obfuscates keys.
