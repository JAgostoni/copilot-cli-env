# Session 2: Environment & Shell Detection Module

## Goal
Build the module responsible for detecting the user's OS, terminal, and shell. This detection must be "advisory" (presented as a default) rather than absolute, as users run complex environments (e.g., WSL, sub-shells).

## Tooling & Library Versions
- **Go Version:** `1.26.x`
- **Shell Detection Libraries:**
  - `github.com/mitchellh/go-ps@v1.0.0` (Recommended for finding the active parent process shell by inspecting `os.Getppid()`).
  - Standard library `os` and `runtime` for environment variables and OS detection.

## Setup Instructions

1. **Install Dependencies:**
   ```bash
   go get -u github.com/mitchellh/go-ps@v1.0.0
   ```

2. **Create Module Structure:**
   ```text
   internal/
     envdetect/
       detect.go
       detect_windows.go
       detect_unix.go
       detect_test.go
   ```

## Tasks

### 1. Define the Detection API
In `internal/envdetect/detect.go`, define the output structure and the main interface:

```go
package envdetect

type Environment struct {
	OS          string // "windows", "darwin", "linux"
	Shell       string // "bash", "zsh", "fish", "powershell", "cmd", "gitbash"
	Terminal    string // "windows-terminal", "alacritty", "vscode", etc. (Optional/Best Effort)
}

// Detector is an interface that allows for mocking in tests
type Detector interface {
	Detect() (Environment, error)
}
```

### 2. Implement OS Detection
Use `runtime.GOOS` to confidently populate the `OS` field.

### 3. Implement Active Shell Detection
Relying purely on `$SHELL` is often inaccurate if a user drops into a sub-shell.

**For Unix (`detect_unix.go`):**
1. Get the parent PID: `ppid := os.Getppid()`
2. Use `go-ps` to find the process name of the parent.
3. Fallback to `os.Getenv("SHELL")` if process traversal fails.
4. Normalize output (e.g., `/bin/zsh` -> `zsh`).

**For Windows (`detect_windows.go`):**
1. Get parent PID and inspect process name (looking for `pwsh.exe`, `powershell.exe`, `cmd.exe`).
2. Fallback to environment heuristics:
   - If `os.Getenv("PSModulePath")` is set, strongly lean toward PowerShell.
   - If `os.Getenv("TERM_PROGRAM")` == `vscode`, it's running inside VS Code.
   - Look for Git Bash specific variables if detecting `bash.exe` on Windows.

### 4. Create the CLI Command Hook
Update `cmd/detect.go` to call your new library:
```go
// Inside detectCmd Run implementation
detector := envdetect.NewDetector()
env, _ := detector.Detect()
fmt.Printf("Detected OS: %s\n", env.OS)
fmt.Printf("Detected Shell: %s\n", env.Shell)
```

## Manual Testing
By the end of this session, perform the following manual tests:
1. Run `go run main.go detect` from `bash` and see "bash".
2. Run `go run main.go detect` from `zsh` and see "zsh".
3. (If on Windows) Run from PowerShell and see "powershell", then run from CMD and see "cmd".
