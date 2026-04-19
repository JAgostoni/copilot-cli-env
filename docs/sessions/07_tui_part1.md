# Session 7: TUI Part 1 - Onboarding & Provider Selection

## Goal
Begin building the interactive terminal UI using the Charm ecosystem to create a guided, keyboard-first onboarding experience.

## Tooling & Library Versions
- **Go Version:** `1.26.x`
- **Libraries:**
  - `charm.land/bubbletea/v2`
  - `github.com/charmbracelet/bubbles@v0.21.1`
  - `charm.land/lipgloss/v2`

## Setup Instructions

1. **Install Dependencies:**
   ```bash
   go get charm.land/bubbletea/v2
   go get github.com/charmbracelet/bubbles@v0.21.1
   go get charm.land/lipgloss/v2
   ```

2. **Structure:**
   ```text
   internal/
     ui/
       app.go          # Main Bubble Tea model
       styles.go       # Lipgloss definitions
       step_env.go     # Step 1: Detect/Confirm Env
       step_provider.go # Step 2: Choose Provider
       step_auth.go    # Step 3: Enter API Key
   ```

## Tasks

### 1. Initialize the Main App Model
In `internal/ui/app.go`, create the root `tea.Model`. It should maintain state about the "current step" and store the accumulated configuration data.

```go
type AppModel struct {
	step        int
	configData  config.CopilotEnv
	// sub-models for different screens
}
```

### 2. Build Step 1: Environment Confirmation
Use the logic from Session 2. Display the detected OS and Shell. Allow the user to press `Enter` to accept, or `e` to edit/override the detected shell using a simple list selector (bubbles/list).

### 3. Build Step 2: Provider Selection
Use `bubbles/list` to present the providers defined in Session 3 (OpenRouter, Anthropic, OpenAI, Ollama). 

### 4. Build Step 3: Authentication Input
Use `bubbles/textinput`.
- If the selected provider requires an API key, show an input field.
- **Crucial:** Set `textinput.EchoMode(textinput.EchoPassword)` to mask the API key.
- If the provider requires a custom base URL (like Azure or Custom), show a second input field for the URL.

## Manual Testing
Run the app using `go run main.go init`. Perform the following manual tests:
- Verify you can navigate from Step 1 (Env) -> Step 2 (Provider) -> Step 3 (Auth) seamlessly using the keyboard.
- Verify that exiting the app at any step cleanly restores the terminal state.
