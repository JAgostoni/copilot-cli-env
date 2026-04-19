# Session 9: Integration, Polish & Security Review

## Goal
Tie the TUI, CLI flags, and Core Logic together into a cohesive binary, handling edge cases and security warnings.

## Tooling & Library Versions
- **Go Version:** `1.26.x`
- No new libraries.

## Tasks

### 1. Connect CLI Flags to TUI Bypass
In `cmd/init.go` and `cmd/configure.go`, write the logic that decides whether to launch the TUI.
- If `go run main.go init` is run with NO flags, launch the full `internal/ui` App.
- If `go run main.go configure --provider openai ...` is run, bypass the TUI entirely. Validate the flags, run the mapping logic, and immediately execute the output renderer.

### 2. Implement the Apply Logic
Wire Step 6 of the TUI (and the `configure` command) to the `internal/profile` manager (Session 6) and `internal/render` (Session 5).
- If the user chose "Profile Update", call `UpdateProfile`.
- If `.env`, write to `./.env`.
- If `Console`, print to `os.Stdout`.

### 3. Security Warning Prompts
Implement a check: if the output mode is `.env` or `Profile`, and the provider has an API key, print a highly visible Lipgloss-styled warning:
`WARNING: You are about to write a plaintext API key to disk.`
Require an extra confirmation step.

### 4. Windows Edge Case Review
Verify that if the OS is Windows:
- Git Bash generates `export` syntax but writes to `~/.bash_profile`.
- PowerShell generates `$env:` syntax.
- Ensure we are not attempting to mutate the Windows Registry.

### 5. Final Build
Create a `Makefile` or simple build script to compile cross-platform binaries:
```bash
GOOS=linux GOARCH=amd64 go build -o bin/copilot-cli-env-linux
GOOS=windows GOARCH=amd64 go build -o bin/copilot-cli-env.exe
GOOS=darwin GOARCH=arm64 go build -o bin/copilot-cli-env-mac
```

## Manual Testing
Perform end-to-end manual testing of the compiled binary:
- Test bypassing the TUI using flags (e.g., `configure --provider openai ...`).
- Test running the tool in different mocked environments or shells (e.g., within Git Bash, VS Code, or standard PowerShell).
- Ensure the output exactly matches GitHub Copilot CLI's BYOK requirements.
