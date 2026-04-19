# Session 1: Foundation & CLI Surface

## Goal
Initialize the project, set up the CLI framework, and define the core data structures that will represent the ultimate environment output required by GitHub Copilot CLI.

## Tooling & Library Versions
- **Go Version:** `1.26.x` (Targeting the latest stable version, 1.26.2, for best cross-platform support).
- **CLI Framework:** `github.com/spf13/cobra@v1.10.2` (The industry standard for Go CLI applications, providing excellent routing, flag parsing, and help generation).

## Setup Instructions

1. **Initialize the Go Module:**
   ```bash
   go mod init github.com/yourusername/copilot-cli-env
   ```

2. **Install Dependencies:**
   ```bash
   go get -u github.com/spf13/cobra@v1.10.2
   ```

3. **Establish Project Structure:**
   Create a standard Go project layout:
   ```text
   cmd/
     root.go        # Cobra root command
     init.go        # Interactive onboarding
     configure.go   # Non-interactive configuration
     detect.go      # Environment detection test command
     models.go      # Model discovery command
     render.go      # Output generation (dry run)
     apply.go       # Persistent profile application
   internal/
     config/        # Core configuration structs
   main.go          # Entrypoint
   ```

## Tasks

### 1. Define Core Configuration Structs
In `internal/config/config.go`, define the data structure that will hold the final Copilot configuration state. This maps directly to the environment variables Copilot expects:

```go
package config

// CopilotEnv represents the final set of environment variables needed.
type CopilotEnv struct {
	ProviderBaseURL string // COPILOT_PROVIDER_BASE_URL (Required)
	Model           string // COPILOT_MODEL (Required)
	ProviderType    string // COPILOT_PROVIDER_TYPE (Optional: 'openai', 'azure', 'anthropic')
	APIKey          string // COPILOT_PROVIDER_API_KEY (Optional for some endpoints)
	Offline         bool   // COPILOT_OFFLINE (Optional: usually 'true' for BYOK)
}
```

### 2. Scaffold CLI Commands
Set up `cobra` commands in the `cmd/` directory.

- `cmd/root.go`: Define the base application description.
- `cmd/init.go`: Define the `init` command (which will eventually launch the TUI).
- Add basic `fmt.Println` placeholders for all commands (`detect`, `models`, `render`, `apply`).

### 3. Implement Basic Flag Parsing
In `cmd/configure.go` (and optionally `root.go` for global flags), set up flags to allow non-interactive usage:
- `--provider`
- `--model`
- `--base-url`
- `--api-key` (Note: ensure we don't print this in default debug output later)
- `--output` (e.g., `env`, `profile`, `console`)

### 4. Wire up `main.go`
Ensure `main.go` simply calls `cmd.Execute()`.

## Manual Testing
By the end of this session, perform the following manual tests:
- Run `go run main.go --help` and see all defined commands.
- Run `go run main.go configure --provider openai --model gpt-4o --output console` and have the command successfully parse the flags (even if it just prints them back for now).
